export default function(context) {

	if(context.route.fullPath == "/login") return

	if( !context.app.$cookies.get("refresh-token") ) {
		context.redirect("/login")
		return
	}

}
 