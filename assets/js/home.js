/**
* Template Name: Siimple - v2.0.1
* Template URL: https://bootstrapmade.com/free-bootstrap-landing-page/
* Author: BootstrapMade.com
* License: https://bootstrapmade.com/license/
*/

// New Vue instance
var App = new Vue({
// Vue instance options here
el: '#app', //
data : {
    adminTeams : null,
    userTeams : null,
    allTeams : null,
    workouts : null,
    detailTeam: null,
},
methods : {
  loadUserData : function(){
    axios({
    method: 'get',
    url: '/api/users/me',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.adminTeams = response.data.AdminTeams
      this.userTeams = response.data.Teams
      this.workouts = response.data.Workouts
      this.detailTeam = response.data.Teams[0]
      })
    .catch(error => {});
  },
  loadAllTeams : function(){
    axios({
    method: 'get',
    url: '/api/teams',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.allTeams = response.data
      })
    .catch(error => {});
  },
  selectTeam : function(team){
    axios({
    method: 'get',
    url: '/api/team/' + str(team.ID) + '/highlight',
    headers: {'Content-Type': 'application/json', 'x-access-token': window.localStorage.user_token}
    })
    .then(response => {
      this.allTeams = response.data
      })
    .catch(error => {});
    $('#teamModal').modal('show')
    this.detailTeam = team
    console.log(team.name)
  },
  joinTeam : function(teamID) {
    console.log(teamID)
    }
  },
})
Vue.config.devtools = true;
App.loadUserData();
App.loadAllTeams();
