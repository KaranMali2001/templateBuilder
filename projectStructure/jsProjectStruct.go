package projectstructure

import (
	"fmt"
	"os"
)

func JsProjectStrcuture(name string) error {
	packageJson := fmt.Sprintf(`
	{
		"name": "%s",
		"version": "1.0.0",
		"scripts": {
			"test": "echo \"Error: no test specified\" && exit 1",
			"dev": "nodemon src/server.js"
		},
		"keywords": [],
		"author": "",
		"license": "ISC",
		"description": "",
		"dependencies": {
			"cors": "^2.8.5",
			"express": "^4.21.0",
			"rimraf": "^6.0.1"
		},
		"devDependencies": {
			"nodemon": "^3.1.7"
		}
	}
	`, name)

	err := os.WriteFile(fmt.Sprintf("%s/package.json", name), []byte(packageJson), 0644)
	if err != nil {
		return err
	}

	directories := []string{
		name + "/src/controllers",
		name + "/src/routes",
		name + "/src/middleware",
		name + "/src/utils",
		name + "/src/config",
	}

	files := map[string]string{
		name + "/src/app.js": `const express = require('express');
const app = express();
const router = require('./routes/index'); 

app.use(express.json());
app.use('/api', router);

module.exports = app;`,

		name + "/src/server.js": `const app = require('./app');
const PORT = 3000;
app.listen(PORT, () => console.log('Server running on port ' + PORT));`,

		name + "/src/routes/index.js": `const express = require('express');
const router = express.Router();
router.get('/', (req, res) => res.send('Hello World!'));
module.exports = router;`,

		name + "/src/controllers/userController.js": `exports.getUser = (req, res) => {
	res.send('Get User');
};`,

		name + "/src/middleware/middleware.js": `// Example authentication middleware
exports.isAuthenticated = (req, res, next) => {
  // Middleware logic
  console.log("middleware called");
  next();
};`,

		name + "/src/utils/utils.js": "",

		name + "/.env": `DATABASE_URL=`,
		name + "/.gitignore": `node_modules,
		
.env
		`,

		name + "/src/config/db.js": `import { PrismaClient } from '@prisma/client';
			const prisma = new PrismaClient();
	export default prisma;`,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("error while creating directory %s: %v", dir, err)
		}
		fmt.Println("Directory created:", dir)
	}

	for filepath, content := range files {
		err := CreateFileWithContent(filepath, content)
		if err != nil {
			return fmt.Errorf("error while creating file %s: %v", filepath, err)
		}
		fmt.Println("File created:", filepath)
	}

	return nil
}
