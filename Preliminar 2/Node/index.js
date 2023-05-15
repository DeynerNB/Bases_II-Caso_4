const express = require('express');
const kafka = require('kafka-node');

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use(express.text());

// Kafka declarations
const client = new kafka.KafkaClient({ kafkaHost: '10.0.0.26:9092' });
const producer = new kafka.Producer(client);

producer.on('ready', () => {
    console.log('Kafka producer is ready');
});
producer.on('error', (err) => {
    console.error('Error with Kafka producer:', err);
});

const topics = [
        {
            topic: 'nuevo_Servicio',
            partitions: 1,
            replicationFactor: 1
        },
        {
            topic: 'actualizar_Servicio',
            partitions: 1,
            replicationFactor: 1
        }
    ];

// API main interface
app.get('/', (req, res) => {

    client.createTopics(topics, (error, result) => {
        if (error) {
            console.error("Topics no creados!")
            res.send("Topics no creados!")
        }
        else {
            console.log("Topics creados exitosamente!");
            res.send("Topics creados exitosamente!")
        }
    });
});

// API add new service
app.post('/agregarServicio', (req, res) => {
    const message = JSON.stringify(req.body);

    const payloads = [
        { topic: 'nuevo_Servicio', messages: message }
    ];

    producer.send(payloads, (err, data) => {
        if (err) {
            console.error('Failed to produce message:', err);
            res.status(500).send('Failed to produce message');
        } else {
            console.log('Message sent:', data);
            res.send('POST request message sent(New data)');
        }
    });
});

// API update specific service
app.post('/actualizarServicio', (req, res) => {
    const message = JSON.stringify(req.body);

    const payloads = [
        { topic: 'actualizar_Servicio', messages: message }
    ];

    producer.send(payloads, (err, data) => {
        if (err) {
            console.error('Failed to produce message:', err);
            res.status(500).send('Failed to produce message');
        } else {
            console.log('Message sent:', data);
            res.send('POST request message sent(Updated data)');
        }
    });
});

app.listen(port, () => {
    console.log(`Server running on port ${port}`);
});