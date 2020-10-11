var app = new Vue({
    el: '#app',
    data() {
        return {
            message: 'Привет, Vue!',
            current_game: {players: [], tasks: []},
            tasks: [],
            labels: [],
            current_task_idx: -1,
            gameSN: '',
            editMode: false,
            adminMode: false,
            game_name: '',
            game_description: ''
        };
    },
    computed: {}

})

Vue.component('task-item', {
    props: ['task', 'idx'],
    methods: {
        onClick: function () {
            changeActiveTask(this.idx);
        },
    },
    template: `
<li>
    <div style="height: 20px">Title: {{ task.title }} x: {{ task.coords.x}} y: {{ task.coords.y}}</div>
</li>
`
})

Vue.component('task-item-editable', {
    props: ['task', 'idx'],
    data: function () {
        return {
            title: '',
            description: ''
        }
    },
    methods: {
        onClick: function () {
            changeActiveTask(this.idx);
        },
        remove: function () {
            removeTask(this.task.idx)
        }
    },
    computed: {},
    template: `
<li>
    <div style="height: 20px" v-on:click="onClick">Title: {{ task.title }} x: {{ task.coords.x}} y: {{ task.coords.y}}</div>
    <input v-model="task.title" placeholder="title">
    <input v-model="task.description" placeholder="description">
    <button v-on:click="this.remove">X</button>
</li>
`
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
    app.current_task_idx = app.tasks.length
}

changeTask = function (idx, task) {
    app.set(app.tasks, idx, task)
}

changeActiveTask = function (idx) {
    app.current_task_idx = idx
}

removeTask = function (idx) {
    app.tasks.splice(app.tasks, idx)
}
removeLastTask = function () {
    app.tasks.pop()
    app.current_task_idx = app.tasks.length
    group.tasks.pop()
}
var feft = function (evt) {
    var coord = map.screenToGeo(evt.currentPointer.viewportX,
        evt.currentPointer.viewportY);
    let x = Math.abs(coord.lat.toFixed(8));
    let y = Math.abs(coord.lng.toFixed(8));
    var tsk = {title: "", coords: {x: x, y: y}}

    tsk.html = function () {
        return `<h6>${tsk.title}</h6>: <p>${tsk.description}</p>`
    }

    addTask(tsk)

    addInfoBubble(map, x, y, tsk)
}

function setUpClickListener(map) {
    // Attach an event listener to map display
    // obtain the coordinates and display in an alert box.
    map.addEventListener('tap', feft);
}

setUpClickListener(map);


createGame = function () {
    axios.post('http://localhost:8080/api/game', {
        name: app.game_name,
        description: app.game_description
    }).then(function (response) {
        app.gameSN = response.data.short_name;
        console.log(response);
        app.editMode = true
        app.adminMode = false
    }).catch(function (error) {
        console.log(error)
    })
}

adminMode = function () {
    axios.get('http://localhost:8080/api/game/' + app.gameSN + '/info')
        .then(function (response) {
            app.current_game = response.data;
            app.game_description = app.current_game.game.description
            app.game_name = app.current_game.game.name
            app.adminMode = true;
            app.editMode = false;
            map.removeEventListener('tap', feft)

            for (let i = 0; i< app.current_game.tasks.length; i++){
                var t = app.current_game.tasks[i]
                addInfoBubble(map, t.coords.x, t.coords.y, t)
            }

            console.log(response);
        })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .then(function () {
            // always executed
        });
}

startQuest = function () {
    axios.post('http://localhost:8080/api/game/' + app.gameSN + '/tasks', {
        tasks: app.tasks
    }).then(function (response) {
        adminMode()
    }).catch(function (error) {
        console.log(error)
    })
}

function addMarkerToGroup(group, coordinate, data) {
    var marker = new H.map.Marker(coordinate);
    // add custom data to the marker
    marker.setData(data);
    group.addObject(marker);
}

function addInfoBubble(map, lat, lng, data) {
    // add 'tap' event listener, that opens info bubble, to the group
    group.addEventListener('tap', function (evt) {
        // event target is the marker itself, group is a parent event target
        // for all objects that it contains
        var bubble = new H.ui.InfoBubble(evt.target.getGeometry(), {
            // read custom data
            content: evt.target.getData().html()
        });
        // show info bubble
        ui.addBubble(bubble);
    }, false);

    addMarkerToGroup(group, {lat: lat, lng: lng}, data);

}

