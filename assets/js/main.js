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
    testValue: "hello world",
    username : null,
    first_name : null,
    last_name : null,
    email : null,
    password : null,
    new_username : null,
    new_password : null,
    confirm_new_password: null,
    loginResponse: null,
    registerResponse: null,
    teams : [],
    stravaID : window.strava_id,
},
methods : {
  registerUser : function(){
    if (this.new_password != this.confirm_new_password){
        this.registerResponse = "Passwords Do Not Match"
        return
    };
    if ((this.email == null) || (this.new_username == null) || (this.new_password == null)){
          this.registerResponse = "All Fields Required"
          return
    };
    var postForm = new FormData();
    postForm.set("username", this.new_username)
    postForm.set("password", this.new_password)
    postForm.set("email", this.email)
    postForm.set("strava_id", this.stravaID)
    axios({
    method: 'post',
    url: '/register',
    data: postForm,
    headers: {'Content-Type': 'multipart/form-data' }
    })
    .then(response => {
      if (response.data['status'] != 200) {
        this.registerResponse = response.data['message']
        return
      }
      token = response.data["token"];
      console.log(token)
      localStorage.setItem('user_token', token)
      window.location.assign("/home");
    })
    .catch(error => {});
  },

  loginFast : function() {
    console.log(this.$strava_id)
    // if (localStorage.user_token != null){
    //   window.location.assign("/home");
    // }
  },

  login : function(){
    var postForm = new FormData();
    postForm.set("username", this.username)
    postForm.set("password", this.password)
    axios({
    method: 'post',
    url: '/login',
    data: postForm,
    headers: {'Content-Type': 'multipart/form-data' }
    })
    .then(response => {
      if (response.data['status'] != 200) {
        this.loginResponse = response.data['message']
        return
      }
      token = response.data["token"];
      console.log(token)
      localStorage.setItem('user_token', token)
      window.location.assign("/home");
    })
    .catch(error => {});
  },
  // login : function(){
  //   if (localStorage.getItem('user-token')) != null {
  //
  //   }
  // },
}


})
Vue.config.devtools = true;


!(function($) {
  "use strict";

  // Toggle nav menu
  $(document).on('click', '.nav-toggle', function(e) {
    $('.nav-menu').toggleClass('nav-menu-active');
    $('.nav-toggle').toggleClass('nav-toggle-active');
    $('.nav-toggle i').toggleClass('bx-x bx-menu');

  });

  // Toogle nav menu drop-down items
  $(document).on('click', '.nav-menu .drop-down > a', function(e) {
    e.preventDefault();
    $(this).next().slideToggle(300);
    $(this).parent().toggleClass('active');
  });

  // Smooth scroll for the navigation menu and links with .scrollto classes
  $(document).on('click', '.nav-menu a, .scrollto', function(e) {
    if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
      e.preventDefault();
      var target = $(this.hash);
      if (target.length) {

        var scrollto = target.offset().top;

        if ($(this).attr("href") == '#header') {
          scrollto = 0;
        }

        $('html, body').animate({
          scrollTop: scrollto
        }, 1500, 'easeInOutExpo');

        if ($(this).parents('.nav-menu').length) {
          $('.nav-menu .active').removeClass('active');
          $(this).closest('li').addClass('active');
          $('.nav-menu').removeClass('nav-menu-active');
          $('.nav-toggle').removeClass('nav-toggle-active');
          $('.nav-toggle i').toggleClass('bx-x bx-menu');
        }
        return false;
      }
    }
  });

})(jQuery);
