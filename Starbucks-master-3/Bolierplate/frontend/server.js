var express = require('express');
var http = require('http');
var path = require('path');
var parser = require('body-parser');
var session = require('client-sessions');
http.post = require('http-post');


app.listen(8000);
console.log('8000 is the magic port');

app.get('/order', function (req, res) {
    console.log("GET /order user:", req.session.user);
    let order = "5a330e54-6236-4677-8675-e7f98c83d863";
    http.get('http://localhost:8080/starbucks/order/' + order, function (response) {
        //console.log("--------" + response);

        response.on('data', function (chunk) {
            console.log(JSON.parse(chunk));
            res.render('pages/order', {
                order: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

