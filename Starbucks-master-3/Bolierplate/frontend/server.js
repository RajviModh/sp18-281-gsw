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

app.post('/order', (request, response) => {
    // let cart = JSON.parse(request.body);
    let username = "";
    if (request.session.user == undefined) {
        username = sessioninfo;
    } else {
        username = req.session.user;
    }

    let items = request.body.items;
    console.log("body:");
    console.log(username);
    console.log(items);
    http.post('http://localhost:8080/starbucks/order', {
        "username": username,
        "items": JSON.stringify(items),
        "location": "San Jose"
    }, (response1) => {
        console.log("statusCode: ", response1.statusCode); // <======= Here's the status code
        response1.on('data', function (chunk) {
            let order = JSON.parse(chunk);

            response.render('pages/order', {
                order: order
            });
        });


    }).on('error', function (e) {
        response.sendStatus(500);
    }).end();

});

app.post('/order/pay', (request, response) => {
    let id = request.body.id;
    let url = "http://localhost:8080/starbucks/order/" + id;

    fetch(url, {
        method: 'PUT', // or 'PUT'
        headers: new Headers({
            'Content-Type': 'application/json'
        })
    }).then(res => res.json()
    )
        .then(jsonResponse => {
            console.log(jsonResponse);
            response.sendStatus(200).end();
        })
        .catch(error => {
            console.error('Error:', error);
            response.sendStatus(400).end();
        });

});

app.get('/getOrders', function (req, res) {
    console.log("kkk" + req.session.user)
    http.get('http://localhost:8080/starbucks/orders/' + req.session.user, function (response) {
        console.log("--------" + response);
        response.on('data', function (chunk) {
            res.render('pages/getOrders', {
                order: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.post('/deleteOrder', function (req, res) {
    var id = req.body.id;
    var abc = encodeURI(id);
    console.log("************************", abc);
    http.post('http://localhost:8080/starbucks/delOrder', {'id': abc}, function (response) {
        response.on('data', function (chunk) {

        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.get('/order/:id', function (req, res) {
    let order = req.params.id;
    console.log("GET /order user:", req.session.user);
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
