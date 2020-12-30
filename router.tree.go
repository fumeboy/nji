package nji

import (
	"net/url"
	"strings"
)

// 路由中的路径参数
type PathParam struct {
	Key   string
	Value string
}

type PathParams []PathParam

// 获得路径参数
func (ps PathParams) Get(name string) (string, bool) {
	for k := range ps {
		if ps[k].Key == name {
			return ps[k].Value, true
		}
	}
	return "", false
}

// 获得路径参数值
func (ps PathParams) Value(name string) (value string) {
	value, _ = ps.Get(name)
	return
}

// 方法树
type methodTree struct {
	method string   // 方法
	root   *node    // 根节点
}

type methodTrees []methodTree

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

type nodeType uint8

const (
	static nodeType = iota
	root
	param
	catchAll
)

type node struct {
	path      string
	// 表示当前节点的path，比如 s，earch，upport 这些字段
	indices   string
	// 通常情况下维护了children列表的path的各首字符组成的string，之所以是通常情况，是在处理包含通配符的path处理中会有一些例外情况
	children  []*node
	handlers  HandlersChain
	// 如果此节点为终结节点handlers为对应的处理链，否则为nil
	priority  uint32
	// 代表了有几条路由会经过此节点，用于在节点进行排序时使用
	nType     nodeType
	// 是节点的类型，默认是static类型，还包括了root类型
	// 对于path包含冒号通配符的情况，nType是 param类型
	// 对于包含 * 通配符的情况，nType类型是 catchAll 类型
	maxParams uint8
	// 是当前节点到各个叶子节点的包含的通配符的最大数量
	wildChild bool
	// 默认是false，当children是 通配符类型时，wildChild为true
	fullPath  string
	// 是从root节点到当前节点的全部path部分
}

// increments priority of the given child and reorders if necessary.
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
			n.indices[pos:pos+1] + // the index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // rest without char at 'pos'
	}

	return newPos
}

func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++
	//计算当前path的通配符的数量
	numParams := countParams(path)

	// 如果单前树是空树，直接在当前node插入path
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(numParams, path, fullPath, handlers)
		n.nType = root
		return
	}
	//path的共同前缀位置
	parentFullPathIndex := 0

walk:
	for {
		// Update maxParams of the current node
		if numParams > n.maxParams {
			n.maxParams = numParams
		}

		// 最长公有前缀
		i := longestCommonPrefix(path, n.path)

		// 如果path与当前的node有部分匹配，需要拆分当前的node
		if i < len(n.path) {
			child := node{
				path:      n.path[i:],  //新的节点包括 node path 没有匹配上的后半部分
				wildChild: n.wildChild,  //设置与当前节点相同
				indices:   n.indices,
				children:  n.children,
				handlers:  n.handlers,
				priority:  n.priority - 1,
				fullPath:  n.fullPath,
			}

			// Update maxParams (max of all children)
			for _, v := range child.children {
				if v.maxParams > child.maxParams {
					child.maxParams = v.maxParams
				}
			}

			n.children = []*node{&child}  //将后半部分设置为孩子节点
			// []byte for proper unicode char conversion, see #65
			n.indices = string([]byte{n.path[i]})
			n.path = path[:i] //当前节点的path只保持前半部分
			n.handlers = nil
			n.wildChild = false // 拆分后的新父节点一定不包含通配符
			n.fullPath = fullPath[:parentFullPathIndex+i] //当前节点的fullPath截取
		}

		// path没有完成匹配，需要继续向下寻找
		if i < len(path) {
			path = path[i:]  //path 更新为没有匹配上的后半部分

			//如果当前节点 wildChild 为 true，那么子节点 children[0] 是通配符节点
			if n.wildChild {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++

				// Update maxParams of the child node
				if numParams > n.maxParams {
					n.maxParams = numParams
				}
				numParams--

				// path为通配符的时候必须一致，然后继续向后
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] {
					// check for longer wildcard, e.g. :name and :names
					// path与当前node的path长度相同 或者path有下划线，继续
					if len(n.path) >= len(path) || path[len(n.path)] == '/' {
						continue walk
					}
				}

				//其他情况会panic
				pathSeg := path
				if n.nType != catchAll {
					pathSeg = strings.SplitN(path, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				panic("'" + pathSeg +
					"' in new path '" + fullPath +
					"' conflicts with existing wildcard '" + n.path +
					"' in existing prefix '" + prefix +
					"'")
			}

			//当前节点的孩子节点不是通配符类型，取出第一个字符
			c := path[0]

			// slash after param
			//冒号通配符后面的 下划线处理
			if n.nType == param && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]  //更新node节点为孩子节点，继续查找
				n.priority++
				continue walk
			}

			// Check if a child with the next path byte exists
			//当前节点的某个孩子与path有相同的前缀
			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]  //更新当前节点为对应的孩子节点，继续查找
					continue walk
				}
			}

			// Otherwise insert it
			//如果是其他情况，新增一个child节点，并且基于这个child节点，插入剩下的path
			if c != ':' && c != '*' {
				// []byte for proper unicode char conversion, see #65
				n.indices += string([]byte{c})
				child := &node{
					maxParams: numParams,
					fullPath:  fullPath,
				}
				n.children = append(n.children, child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			}
			n.insertChild(numParams, path, fullPath, handlers)
			return
		}

		// Otherwise and handle to current node
		if n.handlers != nil {
			panic("handlers are already registered for path '" + fullPath + "'")
		}
		n.handlers = handlers
		return
	}
}


func findWildcard(path string) (wildcard string, i int, valid bool) {
	for start, c := range []byte(path) {
		//如果没有遇到通配符就继续向后查找
		if c != ':' && c != '*' {
			continue
		}

		//找到通配符设置valid为true，那么通配符在path的起始位置就是start
		valid = true
		//从通配符后面继续查找
		for end, c := range []byte(path[start+1:]) {
			switch c {
			//如果遇到下划线，返回wildCard（不包括下划线）、start、true
			case '/':
				return path[start : start+1+end], start, valid

				//如果遇到通配符，valid设置为false
			case ':', '*':
				valid = false
			}
		}
		//在这个位置返回，遍历完了path，valid为true和false的可能性都有
		return path[start:], start, valid
	}
	//在path里没有找到通配符
	return "", -1, false
}



func (n *node) insertChild(numParams uint8, path string, fullPath string, handlers HandlersChain) {
	for numParams > 0 {
		// Find prefix until first wildcard
		wildcard, i, valid := findWildcard(path)
		//path中不包含通配符，直接结束对numParams条件的for循环
		if i < 0 { // No wildcard found
			break
		}

		// valid为false的两种情况是没有找到通配符（之前已经break） 或者 一个path段有多个通配符
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// 如果path段只有通配符没有名字 也会panic；由于wildCard一定是以通配符开头的，通配符后面不能直接加下划线
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		// Check if this node has existing children which would be
		// unreachable if we insert the wildcard here
		// 如果要在当前节点增加一个通配符的孩子节点，当前节点不能有孩子节点，这个时候会导致路由冲突
		if len(n.children) > 0 {
			panic("wildcard segment '" + wildcard +
				"' conflicts with existing children in path '" + fullPath + "'")
		}
		//冒号类型的通配符处理
		if wildcard[0] == ':' { // param
			if i > 0 {
				// Insert prefix before the current wildcard
				n.path = path[:i]  //设置当前节点的path
				path = path[i:]  //更新path
			}
			//孩子节点是通配符，当前节点设置为true
			n.wildChild = true
			child := &node{
				nType:     param,  //冒号类型的通配符类型
				path:      wildcard,  //设置path为wildCard path包含通配符和名字
				maxParams: numParams,
				fullPath:  fullPath,
			}
			n.children = []*node{child}   //children挂接到当前节点
			n = child   //n更新为下沉到孩子节点
			n.priority++
			numParams--  //控制循环的通配符数量减1

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			//如果wildCard的长度小于path，则说明path中还包含以及path
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]  //重新更新path

				child := &node{              //new一个child节点
					maxParams: numParams,
					priority:  1,
					fullPath:  fullPath,
				}
				n.children = []*node{child}
				n = child   //更新n节点为child节点
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			n.handlers = handlers
			return
		}

		//星号通配符类型的处理
		//星号通配符必须是path的最后一个通配符 否则会panic
		if i+len(wildcard) != len(path) || numParams > 1 {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		}

		//星号通配符的前一个字符，必须为下划线，否则panic
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		//当前node的path为星号通配符之前的path
		n.path = path[:i]

		// First node: catchAll node with empty path
		child := &node{  //一个path为空的节点
			wildChild: true,   //空节点的wildCard为true
			nType:     catchAll,
			maxParams: 1,
			fullPath:  fullPath,
		}
		// update maxParams of the parent node
		if n.maxParams < 1 {
			n.maxParams = 1
		}
		n.children = []*node{child}  //空节点挂接到当前node节点
		n.indices = string('/')  //node节点 indices设置为 下划线
		n = child  //node节点下沉为path为空的节点
		n.priority++

		// second node: node holding the variable
		child = &node{
			path:      path[i:],  //path 为从 下划线开始的包含星号通配符的path
			nType:     catchAll,
			maxParams: 1,
			handlers:  handlers,
			priority:  1,
			fullPath:  fullPath,
		}
		n.children = []*node{child}  //将 包含星号通配符的path节点挂接到空节点下

		return
	}

	// 剩下的path不再包含冒号或者星号
	n.path = path
	n.handlers = handlers
	n.fullPath = fullPath
}


// nodeValue holds return values of (*Node).getValue method
type nodeValue struct {
	handlers HandlersChain
	params   PathParams
	tsr      bool
	fullPath string
}

// getValue returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
// nolint:gocyclo
func (n *node) getValue(path string, po PathParams, unescape bool) (value nodeValue) {
	value.params = po
walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if value.handlers = n.handlers; value.handlers != nil {
				value.fullPath = n.fullPath
				return
			}

			if path == "/" && n.wildChild && n.nType != root {
				value.tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			indices := n.indices
			for i, max := 0, len(indices); i < max; i++ {
				if indices[i] == '/' {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == catchAll && n.children[0].handlers != nil)
					return
				}
			}

			return
		}

		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			path = path[len(prefix):]
			// If this node does not have a wildcard (param or catchAll)
			// child,  we can just look up the Next child node and continue
			// to walk down the tree
			if !n.wildChild {
				c := path[0]
				indices := n.indices
				for i, max := 0, len(indices); i < max; i++ {
					if c == indices[i] {
						n = n.children[i]
						continue walk
					}
				}

				// Nothing found.
				// We can recommend to redirect to the same URL without a
				// trailing slash if a leaf exists for that path.
				value.tsr = path == "/" && n.handlers != nil
				return
			}

			// handle wildcard child
			n = n.children[0]
			switch n.nType {
			case param:
				// find param end (either '/' or path end)
				end := 0
				for end < len(path) && path[end] != '/' {
					end++
				}

				// save param value
				if cap(value.params) < int(n.maxParams) {
					value.params = make(PathParams, 0, n.maxParams)
				}
				i := len(value.params)
				value.params = value.params[:i+1] // expand slice within preallocated capacity
				value.params[i].Key = n.path[1:]
				val := path[:end]
				if unescape {
					var err error
					if value.params[i].Value, err = url.QueryUnescape(val); err != nil {
						value.params[i].Value = val // fallback, in case of error
					}
				} else {
					value.params[i].Value = val
				}

				// we need to go deeper!
				if end < len(path) {
					if len(n.children) > 0 {
						path = path[end:]
						n = n.children[0]
						continue walk
					}

					// ... but we can't
					value.tsr = len(path) == end+1
					return
				}

				if value.handlers = n.handlers; value.handlers != nil {
					value.fullPath = n.fullPath
					return
				}
				if len(n.children) == 1 {
					// No handle found. Check if a handle for this path + a
					// trailing slash exists for TSR recommendation
					n = n.children[0]
					value.tsr = n.path == "/" && n.handlers != nil
				}
				return

			case catchAll:
				// save param value
				if cap(value.params) < int(n.maxParams) {
					value.params = make(PathParams, 0, n.maxParams)
				}
				i := len(value.params)
				value.params = value.params[:i+1] // expand slice within preallocated capacity
				value.params[i].Key = n.path[2:]
				if unescape {
					var err error
					if value.params[i].Value, err = url.QueryUnescape(path); err != nil {
						value.params[i].Value = path // fallback, in case of error
					}
				} else {
					value.params[i].Value = path
				}

				value.handlers = n.handlers
				value.fullPath = n.fullPath
				return

			default:
				panic("invalid node type")
			}
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		value.tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.handlers != nil)
		return
	}
}
