import express from 'express';
import { logger } from '../logger.js';
import { upload } from '../middlewares/multer.js';
import config from 'config';
import AdmZip from 'adm-zip';
import fs from 'fs-extra';

const router = express.Router();

router.use('/:page/upload', upload.single('content'));

router.post('/:page/upload', async (req, res) => {
    if (!req.file) {
        res.status(400).send('Bad request');
        return;
    }

    const site = config.get('sites')[req.params.page];
    if (!site) {
        res.status(404).send('Not found');
        return;
    }

    logger.info(`Received upload for '${req.params.page}'`);
    const decompressed = `/tmp/www-uploads/${req.file.filename}-decompressed`;
    
    let admZip: AdmZip;
    try {
        admZip = new AdmZip(req.file.path);
    } catch (e) {
        logger.error(`Error creating zip file: ${e}`);
        res.status(500).send('Internal server error');
        return;
    }

    logger.info(`Decompressing file to ${decompressed}`);
    try {
        admZip.extractAllTo(decompressed);
    } catch (e) {
        logger.error(`Error decompressing file: ${e}`);
        res.status(500).send('Internal server error');
        return; 
    }
    
    logger.info(`Moving decompressed files to ${site.path}`);
    try {
        fs.emptyDirSync(site.path);
        fs.readdirSync(decompressed).forEach(file => {
            fs.moveSync(`${decompressed}/${file}`, `${site.path}/${file}`);
        });
    } catch (e) {
        logger.error(`Error moving files: ${e}`);
        res.status(500).send('Internal server error');
        return; 
    }
    
    logger.info(`Deleting temporary files`);
    try {
        fs.removeSync(decompressed);
        fs.removeSync(req.file.path);
    } catch (e) {
        logger.error(`Error deleting temporary files: ${e}`);
        res.status(500).send('Internal server error');
        return; 
    }
    
    const endDate = new Date();
    logger.info(`Upload complete. New site published at ${endDate.toString()}.`);
    logger.info('');
    
    res.status(200).send('OK');
});

export default router;
