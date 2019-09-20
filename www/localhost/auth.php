<?php
    session_start();
    if (isset($_SESSION['username']) or isset($_SESSION['id'])){
        header('location: index.php');
        die();
    }
?>

<!DOCTYPE html>
<html>
<head>
	<title>Authorization</title>
	<meta charset="utf-8">
</head>
<body>
	<div class="enter">
		<form method="post" action="add_user.php">
			<div class="t">
				<input class="inp" type="text" name="login" placeholder="Логин:">
				<input class="but" type="submit" name="enter" value="Войти">
			</div>
		</form>
	</div>
</body>
</html>