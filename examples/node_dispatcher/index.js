'use strict';

// This example is primarily used to benchmark the Node.js
// dispatcher against golang dispatcher

const async  = require('async');
const Docker = require('dockerode');
const amqp   = require('amqplib/callback_api');

const docker = new Docker({host: 'http://172.17.0.1', port: 2375});
let counter = 0;

amqp.connect('amqp://localhost', function(err, conn) {
  conn.createChannel(function(err, ch) {
    var q = 'myqueue';
    ch.assertQueue(q, {durable: false});
    ch.consume(q, function(msg) {
        const obj = JSON.parse(msg.content.toString());
        docker.run(obj.Image, obj.Cmd, process.stdout, function (err, data, container) {
            // remove after exit
            container.remove(function () {});
            if (++counter === 10) {
                process.exit(0);
            }
        });
    }, {noAck: true});
  });
});
