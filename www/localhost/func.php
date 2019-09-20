<?php
function MakeRequest($dir = "", $query = "")
    {
        $url = "http://server_api:9000/".$dir;
        $myCurl = curl_init();
        curl_setopt_array($myCurl, array(
            CURLOPT_URL => $url,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => $query,
            CURLOPT_HEADER => "Content-Type: application/json"
        ));
        $response = curl_exec($myCurl);
        curl_close($myCurl);
        return $response;
    }
?>