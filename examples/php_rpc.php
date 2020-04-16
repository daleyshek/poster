<?php

# PHP的实现示例
# 示例通过curl实现调用json-rpc服务
# json-rpc规范 http://wiki.geekdream.com/Specification/json-rpc_2.0.html
# 运行前需要启动json-rpc服务，即运行poster可执行文件

function call(){
    $req = <<<json
    {
        "jsonrpc": "2.0",
        "method": "Poster.Generate",
        "params": [
        {
            "title": "BestFriendsChina",
            "content": "凡读书......须要读得字字响亮，不可误一字，不可少一字，不可多一字，不可倒一字，不可牵强暗记，只是要多诵数遍，自然上口，久远不忘。古人云，“读书百遍，其义自见”。谓读得熟，则不待解说，自晓其义也。余尝谓，读书有三到，谓心到，眼到，口到。心不在此，则眼不看仔细，心眼既不专一，却只漫浪诵读，决不能记，记亦不能久也。三到之中，心到最急。心既到矣，眼口岂不到乎？",
            "imageURL": "https://api.mcdsh.com/storage/images/800/g45rDYNgaI3z0ZACqzkI0iysuXIz4omyBhZSGBUM.jpeg",
            "QRCodeURL": "resources/qrcode.jpg"
        }],
        "id": 1
    }
json;

    $ch = curl_init();
    curl_setopt($ch,CURLOPT_URL,"http://127.0.0.1:2019");
    curl_setopt($ch,CURLOPT_RETURNTRANSFER,1);
    curl_setopt($ch,CURLOPT_POSTFIELDS,$req);
    curl_setopt($ch,CURLOPT_POST,true);
    curl_setopt($ch,CURLOPT_HTTPHEADER,["Content-Type: application/json-rpc"]);

    $res = curl_exec($ch);
    if($res === false) {
        echo "fail";
        return;
    }
    curl_close($ch);
    var_dump($res);
    //将返回
    //{"id":1,"result":"ubyFS5l3UmkCG2dd8Gpfj7MkETtGqPFMkMeZkBDSTKIP5xkNJr7djI3gL3gqgo.jpg","error":null}
}

call();