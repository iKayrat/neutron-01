<!DOCTYPE html>
<html lang="en">

<head>
  <title>{{.Title}}</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css"
    integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
  <link rel="stylesheet" href="C:/Users/User/go/src/bootstrap/css/bootstrap.min.css"
    integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <link href="/static/css/main.css" rel="stylesheet">
</head>
<div class="header">
  <a href="/" class="logo">Beego Test</a>
  <div class="header-right">
    <a href="/">Home</a>
    <a href="">Contact</a>
    <a href="">About</a>
    <div class="header-right">
      {{if .InSession}}
      <a href="#">
      {{.Firstname}}
      </a>
      <a href="/user/logout">Logout</a>
      <a href="/profile">Profile</a>
      {{else}}
      <a class="active" href="/auth/login/home">Log In</a>
      <a class="btn btn-success" href="/auth/register"><i class=""></i>Sign Up</a>
      {{end}}
    </div>
  </div>

</div>