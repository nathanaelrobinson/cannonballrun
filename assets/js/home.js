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
    adminTeams : '',
    userTeams : '',
    allTeams : '',
    joinableTeams: '',
    workouts : '',
    detailTeam: '',
    detailRunners: '',
    new_teamname: '',
    new_description: '',
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
  loadJoinableTeams : function(){
    $('#joinTeamModal').modal('show')
    axios({
    method: 'get',
    url: '/api/teams',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.joinableTeams = response.data
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
    this.detailTeam = team
    axios({
    method: 'post',
    url: '/api/teams/' + this.detailTeam.ID + '/join',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      $('#joinTeamModal').modal('hide')
      window.alert('Team Join Pending Admin Approval')
      })
    .catch(error => {});
  },

  convertToDate : function(timestring) {
    unix = Date.parse(timestring)
    s = new Date(unix).toLocaleDateString("en-US")
    return s
  },
  createNewTeam : function(){
    var postForm = new FormData();
    postForm.set("name", this.new_teamname)
    postForm.set("description", this.new_description)
    axios({
    method: 'post',
    url: 'api/teams/0',
    data: postForm,
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
        $('#newTeamModal').modal('hide')
        return
    })
    .catch(error => {});
  },
},
})


Vue.config.devtools = true;
App.loadWorkouts();
App.loadAllTeams();
App.loadMyTeams();
