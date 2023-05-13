const express = require('express');
const kafka = require('kafka-node');

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use(express.text());

// Kafka declarations
const client = new kafka.KafkaClient({ kafkaHost: '10.0.0.51:9092' });
const producer = new kafka.Producer(client);

producer.on('ready', () => {
    console.log('Kafka producer is ready');
});
producer.on('error', (err) => {
    console.error('Error with Kafka producer:', err);
});

// API main interface
app.get('/', (req, res) => {
    res.send("Select another option")
});

// API add new service
app.post('/agregarServicio', (req, res) => {
    const message = req.body;

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
    const message = req.body;

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