<?php
    session_start();
    require("func.php");
    if (!isset($_SESSION['username']) or !isset($_SESSION['id'])){
        header('location: auth.php');
        die();
    }
    if (isset($_POST['send']) and !empty($_POST['text'])) {
        $query = json_encode(array('chat' => $_POST['chat'], 'author' => $_SESSION["id"], 'text' => $_POST["text"]), JSON_PRETTY_PRINT);
        $response = MakeRequest("messages/add", $query);
        $res = json_decode($response, true);
        if (!isset($res)) {
            $_SESSION["error"] = "Incorrect input";
        }
        if ($res["status"] != "ok") {
            $_SESSION["error"] = $res["status"];
        }
    } else {
        $_SESSION["error"] = "Empty input";        
    }
    header('location: messages.php?chat_id='.$_POST['chat']);
    die();
?>