<template>
	<div class="wrapper">
		<c-navbar />
		<c-aside />
		<nuxt />
		<c-footer />
	</div>
</template>
<script>
import Aside	from "./components/aside.vue"
import Navbar	from "./components/navbar.vue"
import Footer	from "./components/footer.vue"
import "~/assets/css/app.css"

export default {
	components: {
		"c-aside": Aside,
		"c-navbar": Navbar,
		"c-footer": Footer
	}
}
</script>
