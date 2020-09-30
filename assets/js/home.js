Vue.component('team-modal', {
    template: '#modal-template',
    props: ['team'],
    data: function() {
        return {

        }
    }
});
// New Vue instance
var App = new Vue({
// Vue instance options here
el: '#app', //
data : {
    adminTeams : null,
    userTeams : null,
    allTeams : null,
    workouts : null,
    detailTeam: '',
    detailRunners: null,
},
methods : {
  loadWorkouts : function(){
    axios({
    method: 'get',
    url: '/api/users/me/workouts',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.workouts = response.data
      })
    .catch(error => {});
  },
  loadAllTeams : function(){
    axios({
    method: 'get',
    url: '/api/teamsstatus',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.allTeams = response.data
      })
    .catch(error => {});
  },
  loadMyTeams : function(){
    axios({
    method: 'get',
    url: '/api/users/me/affiliations',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.userTeams = response.data
      })
    .catch(error => {});
  },
  loadRunners : function(team){
    axios({
    method: 'get',
    url: '/api/team/details/' + team.id,
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.detailRunners = response.data
      return
      })
    .catch(error => {});
  },
  selectTeam : function(team){
    this.detailTeam = team
    this.loadRunners(team)
    // axios({
    // method: 'get',
    // url: '/api/team/' + str(team.ID) + '/highlight',
    // headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    // })
    // .then(response => {
    //   this.allTeams = response.data
    //   })
    // .catch(error => {});
    $('#teamModal').modal('show')
    // this.detailTeam = team
    console.log(team.name)
  },
  joinTeam : function(team) {
    axios({
    method: 'post',
    url: '/api/teams/' + team.ID + '/join',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      console.log(response)
      })
    .catch(error => {});
  },

  convertToDate : function(timestring) {
    unix = Date.parse(timestring)
    s = new Date(unix).toLocaleDateString("en-US")
    return s
  },
},
})


Vue.config.devtools = true;
App.loadWorkouts();
App.loadAllTeams();
App.loadMyTeams();
