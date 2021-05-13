let feed = new Vue({
    el: '#stream',
    delimiters: ["[[", "]]"],
    data: {
	offset: -10,
	limit: 10,
	activities: []
    },
    methods: {
	fetchActivities() {
	    this.offset += 10
	    fetch(`/activities?offset=${this.offset}&limit=${this.limit}`)
	    .then(res => {
		return res.json()
	    })
	    .then(data => {
		console.log("data: ", data)
		this.activities = this.activities.concat(data)
	    })
	    .catch(err => console.log("fetch error: ", err))
	}
    },
    mounted() {
	let vueInstance = this
	let target = document.querySelector('#scroll_target')
	let options = {
	    root: null,
	    threshold: 0.8
	}
	let observer = new IntersectionObserver((entries, o) => {
	    console.log("entered viewport")
	    vueInstance.fetchActivities()
	}, options)
	observer.observe(target)
    }
})
