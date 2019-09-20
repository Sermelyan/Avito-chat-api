<?php
    session_start();
    require("func.php");
    if (isset($_SESSION['username']) or isset($_SESSION['id'])){
        header('location: index.php');
        die();
    }
    if (isset($_POST['enter']) and !empty($_POST['login'])) {
        $query = json_encode(array('username' => $_POST['login']), JSON_PRETTY_PRINT);
        $response = MakeRequest("users/add", $query);
        $res = json_decode($response, true);
        if ($res["status"] == "ok") {
            $_SESSION["username"] = $_POST['login'];
            $_SESSION["id"] = $res["id"];
            header('location: index.php');
            die();
        } 
    }else {
        header('location: auth.php');
        die();
    }

?>