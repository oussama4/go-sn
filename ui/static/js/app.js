let feed = new Vue({
    el: '#stream',
    delimiters: ["[[", "]]"],
    data: {
	offset: 0,
	limit: 10,
	activities: []
    },
    methods: {
	fetchActivities() {
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
    created() {
	this.fetchActivities()
    }
})
