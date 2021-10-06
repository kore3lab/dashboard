<template>
	<div class="wrapper toggled">
		<c-navbar />
		<c-aside />
		<nuxt />
    <c-log/>
		<c-footer />
	</div>
</template>
<script>
import Aside	from "./components/aside.vue"
import Navbar	from "./components/navbar.vue"
import Log	from "./components/log.vue"
import Footer	from "./components/footer.vue"

export default {
	components: {
		"c-aside": Aside,
		"c-navbar": Navbar,
    "c-log" : Log,
		"c-footer": Footer
	}
}
</script>