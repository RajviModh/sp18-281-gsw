var express = require('express');
var http = require('http');
var path = require('path');
var parser = require('body-parser');
var session = require('client-sessions');
http.post = require('http-post');

var app = express();
var sessioninfo;

app.use(session({
    cookieName: 'session',
    secret: 'random_string_goes_here',
    duration: 30 * 60 * 1000,
    activeDuration: 5 * 60 * 1000,
}));

app.use(express.static("Views"));

// set the view engine to ejs
app.set('view engine', 'ejs');
app.use(parser.urlencoded({extended: false}));
app.use(parser.json());
// use res.render to load up an ejs view file

app.get('/', function (req, res) {
    console.log("In GET /");
    /*http.get('http://localhost:8080/', function (response) {
        //console.log(response.body);
        response.on('data', function (chunk) {
            res.render('pages/login', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();*/
    res.render('pages/login');
});

app.get('/signup', function (req, res) {
    http.get('http://localhost:8080/', function (response) {
        //console.log(response.body);
        response.on('data', function (chunk) {
            res.render('pages/signup', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.post('/proceedToCheckout', function (req, res) {
    console.log("Here where you want");
    res.render('pages/getOrders')
});

// Add items to cart request
app.post('/addToCart', function (req, res) {
    //console.log("sssss"+req.session.user)
    //console.log(sessioninfo)
    http.post('http://localhost:8080/starbucks/addToCart', {
        'name': req.body.name,
        'price': req.body.price,
        'calories': req.body.calories,
        'username': sessioninfo
    }, function (response) {
        if (response.statusCode == 200) {
            // console.log("Item added to cart");
        } else {
            console.log("Item could not be added.Please try again..");
        }
        //console.log(response.body);
        response.on('data', function (chunk) {
            res.redirect('pages/home');
            /*res.render('pages/home', {
                resp: JSON.parse(chunk)
            });*/
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.post('/submitSignUp', function (req, res) {
    var fname = req.body.fname;
    var lname = req.body.lname;
    var uname = req.body.username;
    var pass = req.body.cpassword;
    var loc = req.body.location;
    http.post('http://localhost:8080/submitSignUp', {
        'fname': fname,
        'lname': lname,
        'email': uname,
        'password': pass,
        'location': loc
    }, function (response) {
        console.log("statusCode: ", response.statusCode); // <======= Here's the status code
        if (response.statusCode == 200) {
            res.render('pages/login', {});
        }
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.get('/menu', (request, response) => {
    let username = request.session.user;
    http.get('http://localhost:8080/starbucks/getMenu', function (response1) {
        //console.log("--------" + response);
        console.log("After response");
        response1.on('data', function (chunk) {
            response.render('pages/home', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        response.sendStatus(500);
    }).end();
});

app.post('/login', function (req, res) {
    var uname = req.body.username;
    var pass = req.body.password;
    req.session.user = uname;
    sessioninfo = uname;
    http.post('http://localhost:8080/login', {'email': uname, 'password': pass}, function (response) {
        console.log("statusCode: ", response); // <======= Here's the status code
        if (response.statusCode == 200) {
            console.log("kkkk" + req.session.user)
            http.get('http://localhost:8080/starbucks/getMenu', function (response) {
                //console.log("--------" + response);
                console.log("After response");
                response.on('data', function (chunk) {
                    res.render('pages/home', {
                        resp: JSON.parse(chunk)
                    });
                });
            }).on('error', function (e) {
                res.sendStatus(500);
            }).end();
        }
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

/*

app.post('/addToCart', function (req, res) {
    console.log("sssss" + req.session.user)
    http.post('http://localhost:8080/starbucks/addToCart', {
        'name': req.body.name,
        'price': req.body.price,
        'calories': req.body.calories,
        'username': req.session.user
    }, function (response) {
        //console.log(response.body);
        if (response.statusCode == 200) {
            window.console.log("Item added to cart");
        } else {
            window.console.log("Item could not be added.Please try again..");
        }
        response.on('data', function (chunk) {
            res.render('pages/home', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});
*/

app.post('/incrementItem', function (req, res) {
//'{"name": ' + req.body.name + ',"price": ' + req.body.price + ',"calories": ' + req.body.calories + ',"username": ' + req.session.user + '}'

    console.log("POST /incrementItem " + req.session.user);
    console.log(req.body.name);
    console.log(req.body.price);
    console.log(req.body.calories);
    console.log(req.session.user);
    console.log(sessioninfo);


    let name = req.body.name;
    let price = req.body.price;
    let calories = req.body.calories;
    let username = "";
    if (req.session.user == undefined) {
        username = sessioninfo;
    } else {
        username = req.session.user;
    }
    http.post('http://localhost:8080/starbucks/addToCart', {
        'name': name,
        'price': price,
        'calories': calories,
        'username': username
    }, function (response) {
        //console.log(response.body);
        if (response.statusCode == 200) {
            console.log("Item incremented");
        } else {
            console.log("Item could not be incremented.Please try again..");
        }
        response.on('data', function (chunk) {
            console.log("cart:" + chunk);
            res.render('pages/cart', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();


});


app.get('/logout', function (req, res) {
    console.log("sssss" + req.session.user)
    http.get('http://localhost:8080/logout', function (response) {
        req.session.reset();
        console.log("statusCode: ", response); // <======= Here's the status code
        if (response.statusCode == 200) {
            console.log("sssss" + req.session.user);

            res.render('pages/login', {});
        }
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});


/*app.get('/home', function(req, res) {
    http.get('http://localhost:8080/', function(response) {
        console.log("statusCode: ", response); // <======= Here's the status code
        if(response.statusCode == 200) {
            res.render('pages/home', {
            });
        }
    }).on('error', function(e) {
        res.sendStatus(500);
    }).end();
});*/

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

/*app.post('/updateOrder', function(req, res) {
    var id = req.body.id;
    console.log(id);
    var abc = encodeURI(id);
    console.log("************************",abc);
    http.post('http://localhost:8080/starbucks/updateOrder',{'id':abc}, function(response) {
        //console.log("----update order----",response);
        response.on('data',function (chunk){
        });
    }).on('error', function(e) {
        res.sendStatus(500);
    }).end();
});*/

app.get('/cart', function (req, res) {
    console.log("^^^^^^^^^session^^^^^^^^^^^^^^", req.session.user);
    http.get('http://localhost:8080/starbucks/cart/' + req.session.user, function (response) {
        //console.log("--------" + response);

        response.on('data', function (chunk) {
            console.log(JSON.parse(chunk));
            res.render('pages/cart', {
                resp: JSON.parse(chunk)
            });
        });
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});
app.post('/addItem', function (req, res) {
    var id = req.body.id;
    var data = req.body.data
    var abc = encodeURI(id);
    console.log("************************", abc, data);
    http.post('http://localhost:8080/starbucks/cart/addItem', {'id': abc, 'data': data}, function (response) {
        console.log("statusCode: ", response.statusCode); // <======= Here's the status code
        if (response.statusCode == 200) {
            // response.redirect('http://localhost:8000/cart')
            res.render('pages/home', {});
        }
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.post('/deleteItem', function (req, res) {
    var id = req.body.id;
    var data = req.body.data;
    var abc = encodeURI(id);
    console.log("************************", abc, data);
    http.post('http://localhost:8080/starbucks/cart/quantity', {'id': abc, 'data': data}, function (response) {
        console.log("statusCode: ", response.statusCode); // <======= Here's the status code
        if (response.statusCode == 200) {
            res.render('pages/home', {});
        }
    }).on('error', function (e) {
        res.sendStatus(500);
    }).end();
});

app.post('/deleteCart', function (req, res) {
    var id = req.body.id;
    var data = req.body.data;
    var abc = encodeURI(id);
    console.log("************************", abc, data);
    http.post('http://localhost:8080/starbucks/cart/delete', {'id': abc, 'data': data}, function (response) {
        console.log("statusCode: ", response.statusCode); // <======= Here's the status code
        if (response.statusCode == 200) {
            res.render('pages/home', {});
        }
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
            console.log("order:" , JSON.stringify(order));
            response.status(201).json(order);
        });

    }).on('error', function (e) {
        response.sendStatus(400);
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


app.listen(8000);
console.log('8000 is the magic port');
