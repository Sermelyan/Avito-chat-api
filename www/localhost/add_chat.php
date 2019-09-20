<?php
    session_start();
    require("func.php");
    if (!isset($_SESSION['username']) or !isset($_SESSION['id'])){
        header('location: auth.php');
        die();
    }
    if (isset($_POST['create']) and !empty($_POST['name']) and !empty($_POST['users'])) {
        $users = explode(',', $_POST['users']);
        $query = json_encode(array('name' => $_POST['name'], 'users' => $users), JSON_PRETTY_PRINT);
        $response = MakeRequest("chats/add", $query);
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
    header('location: index.php');
    die();
?>