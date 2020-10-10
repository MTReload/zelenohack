var app = new Vue({
    el: '#app',
    data() {
        return {
            message: 'Привет, Vue!',
            current_game: null,
            tasks: [],
            current_task_idx: -1,
        };
    },
    created() {
        axios
            .get('https://mtreload.ru/api/game/abcde/info')
            .then(response => (this.current_game = response.data));
    },
    computed: {}

})


Vue.component('task-item', {
    props: ['task'],
    template: '<span>Title: {{ task.title }} x: {{ task.coords.x}} y: {{ task.coords.y}}</span>>'
})

Vue.component('player-item', {
    props: ['player', 'tasks'],
    computed: {
        getTask: function () {
            for (let i = 0; i < this.tasks.length; i++) { // выведет 0, затем 1, затем 2
                if (this.tasks[i].id === this.player.task_id) {
                    return this.tasks[i]
                }

            }
        }
    },
    template: `<li>{{ player.name }} |
    
    <task-item v-bind:task="getTask"></task-item>

</li>`
})

addTask = function (task) {
    app.tasks.push(task)
}

changeTask = function (idx, task) {
    app.set(app.tasks, idx, task)
}

removeTask = function () {
    app.tasks.pop()
}

function setUpClickListener(map) {
    // Attach an event listener to map display
    // obtain the coordinates and display in an alert box.
    map.addEventListener('tap', function (evt) {
        var coord = map.screenToGeo(evt.currentPointer.viewportX,
            evt.currentPointer.viewportY);
        logEvent('Clicked at ' + Math.abs(coord.lat.toFixed(4)) +
            ((coord.lat > 0) ? 'N' : 'S') +
            ' ' + Math.abs(coord.lng.toFixed(4)) +
            ((coord.lng > 0) ? 'E' : 'W'));
    });
}