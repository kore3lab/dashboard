<template>
<div class="wrapper">
	<c-navbar />
	<c-aside />
    <nuxt />
	<c-footer />
</div>
</template>

<script>
import Aside	from './components/aside.vue'
import Navbar	from './components/navbar.vue'
import Footer	from './components/footer.vue'
// css
import '@fortawesome/fontawesome-free/css/all.css'
import '~/static/admin-lte/css/adminlte.css'
// javascript
import 'expose-loader?$!expose-loader?jQuery!jquery'
import '~/static/admin-lte/js/adminlte'

export default {
	components: {
		"c-aside": Aside,
		"c-navbar": Navbar,
		"c-footer": Footer
	},
	head: {
		script: [
			// { src: "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.min.js" }
		],
		link: [
			// { rel: "stylesheet", href: "https://fonts.googleapis.com/css?family=Roboto" }
		]
	}
}
</script>
