import multer from 'multer';

const storage = multer.diskStorage({
    destination: '/tmp/www-uploads',
    filename: function (req, file, cb) {
        const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9)
        cb(null, file.fieldname + '-' + uniqueSuffix)
      }
});
export const upload = multer({ 
    storage: storage,
    limits: {
        fileSize: 1024 * 1024 * 500 
    },
    fileFilter: function (req, file, cb) {
        if (!file.originalname.match(/\.(zip)$/)) {
            return cb(new Error('Only accepting zip files'));
        }
        cb(null, true);
    }
});
