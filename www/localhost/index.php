<?php 
    require("func.php");
    session_start();
    if (!isset($_SESSION['username']) or !isset($_SESSION['id'])){
		header('location: auth.php');
		die();
	}
	echo '<!DOCTYPE html>
	<html>
	<head>
		<title>Main</title>
		<meta charset="utf-8">
	</head>
	<body>
		<div class="enter">
			<form method="post" action="add_chat.php">
				<div class="t">
					<input class="inp" type="text" name="name" placeholder="Chat name"> 
					Please enter users id separeted by comma without whitespaces. Dont forget your id):
					<input class="inp" type="text" name="users" placeholder="Users ">
					<input class="but" type="submit" name="create" value="Create">
				</div>
			</form>
		</div>';
	echo "Hello user ".$_SESSION['username'];
	echo ". Your id ".$_SESSION['id'];
	if (isset($_SESSION["error"])) {
		echo "<div class='error'>".$_SESSION["error"]."</div>";
		$_SESSION["error"] = NULL;
	}
	$query = json_encode(array("user" => $_SESSION["id"]), JSON_PRETTY_PRINT);
	$response = MakeRequest("chats/get", $query);
	$res = json_decode($response, true);
	echo '<table border="1px" bordercolor="black">';
	echo '<tr><th>id</th><th>name</th><th>users</th><th>Created</th><th>Last message</th></tr>';
	foreach ($res as $_ => $value) {
		echo '<tr>';
		foreach ($value as $key1 => $value1) {
			switch ($key1) {
				case 'id':
					echo '<td><a href="messages.php?chat_id='.$value1.'">Go to chat</a></td>';
					break;
				case 'name':
					echo '<td>'.$value1.'</td>';					
					break;
				case 'users':
					echo "<td>";
					foreach ($value1 as $_ => $user) {
						echo $user." ";
					}
					echo "</td>";
					break;
				case 'last_message':
					if ($value1 == "0001-01-01T00:00:00Z") {
						echo "<td>No messages</td>";
					} else {
						echo "<td>".$value1."</td>";					
					}
					break;
				default:
					echo "<td>".$value1."</td>";
					break;
			}
		}
		echo "</tr>";
	}
	echo "</table>";
	echo '<a href="exit.php">Logout</a>
	</body>
	</html>';
?>
