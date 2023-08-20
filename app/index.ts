import { logger } from './logger.js';
import express from 'express';
import upload from './routes/upload.js';
import config from 'config';

const startDate = new Date();

logger.info(`Welcome to scrapbook, ${startDate.toString()}`);
logger.info('');

const app = express();

app.use(upload);

app.listen(config.get('port'), () => {
    logger.info(`Server started on port ${config.get('port')}`);
});
