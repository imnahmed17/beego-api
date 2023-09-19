<!DOCTYPE html>

<html>
<head>
  <title>Beego</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link rel='stylesheet' href='https://cdn-uicons.flaticon.com/uicons-regular-rounded/css/uicons-regular-rounded.css'>
  <link rel="stylesheet" href="static/css/t-datepicker.min.css">
  <link rel="stylesheet" href="static/css/t-datepicker-main.css">
</head>

<body>
  <header>
    <h1>Welcome to Beego</h1>
    <form action="" method="get">
      <label for="location">Location</label>
      <input type="text" name="location">
      <br>
      <div class="t-datepicker">
        <div class="t-check-in"></div>
        <div class="t-check-out"></div>
      </div>
      <br><br><br>
      <label for="page">Page No</label>
      <input type="number" name="page">
      <input type="submit" value="search">
    </form>
  </header>
  <footer>
    <div class="author">
      Official website:
      <a href="http://{{.Website}}">{{.Website}}</a> /
      Contact me:
      <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
    </div>
  </footer>
  <div class="backdrop"></div>

  <script src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
  <script src="static/js/t-datepicker.min.js"></script>
  <script src="static/js/t-datepicker.js"></script>
</body>
</html>
