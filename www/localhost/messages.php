<?php 
    session_start();
	require("func.php");
    if (!isset($_SESSION['username']) or !isset($_SESSION['id'])){
        header('location: auth.php');
        die();
    }
    $id = $_GET['chat_id'];

    echo '<!DOCTYPE html>
	<html>
	<head>
		<title>Chat</title>
		<meta charset="utf-8">
	</head>
    <body>
        <a href="index.php">To main</a>';

    $query = json_encode(array('chat' => $id), JSON_PRETTY_PRINT);
    $response = MakeRequest("messages/get", $query);
    $res = json_decode($response, true);
    echo '<table border="1px" bordercolor="black">';
	echo '<tr><th>Author</th><th>Text</th><th>Created</th></tr>';
	foreach ($res as $_ => $value) {
		echo '<tr>';
		foreach ($value as $key1 => $value1) {
			switch ($key1) {
				case 'author':
				case 'text':
				case 'created_at':
                    echo "<td>".$value1."</td>";
					break;
			}
		}
		echo "</tr>";
	}
    echo "</table>";
    echo '<div class="enter">
            <form method="post" action="add_mes.php">
                <div class="t">
                    <input class="inp" type="text" name="text" placeholder="Input text"> 
                    <input class="but" type="submit" name="send" value="Send">
                    <input type="hidden" name="chat" value="'.$id.'">
                </div>
            </form>
            <a href="exit.php">Logout</a>
        </div>
	</body>
	</html>';