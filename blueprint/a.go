package blueprint

type ctx struct {

}

func (c *ctx) PathParam(){

}

type view struct {
	A string
}

func (v *view) handle(ctx *ctx){
	v.A = PathParam(IsPhoneNum, Between(1,9))

}


