package projectstructure

import (
	"fmt"
	"os"
	"os/exec"
)

func TsProjectStructure(name string) error {
	packageJson := fmt.Sprintf(`
  {
  "name": "%s",
  "version": "1.0.0",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1",
    "build": "npx tsc",
    "start": " node dist/server.js",
    "dev": " tsc -w & nodemon dist/server.js"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "description": "",
  "dependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.6.1",
    "cors": "^2.8.5",
    "express": "^4.21.0",
    "rimraf": "^6.0.1"
  },
  "devDependencies": {
    "nodemon": "^3.1.7"
  }
}

`, name)
	os.WriteFile(fmt.Sprintf("%s/package.json", name), []byte(packageJson), 0644)

	directories := []string{
		name + "/src/controllers",
		name + "/src/routes",
		name + "/src/middleware",
		name + "/src/utils",

		name + "/src/config",
	}
	files := map[string]string{
		name + "/src/app.ts": `import express from 'express';
const app = express();
import router from './routes/index';  // Assuming you have a router in 'routes/index'

app.use(express.json());
app.use('/api', router);

export default app;`,

		name + "/src/server.ts": `import app from './app';
const PORT = 3000;
app.listen(PORT, () => console.log('Server running '));`,
		name + "/src/routes/index.ts": `import { Router ,Request,Response} from 'express';
		const router = Router();
		router.get('/', (req:Request, res:Response) => res.send('Hello World!'));
		export default router;`,

		name + "/src/controllers/userController.ts": `import { Request, Response } from "express";
		export const getUser = (req:Request, res:Response) => res.send('Get User');`,
		name + "/src/middleware/middleware.ts": `import { Request, Response, NextFunction } from "express";

// Example authentication middleware
export const isAuthenticated = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  //write your middleware logic here
  console.log("middleware called");
  next();
};
`,
		name + "/src/utils/utils.ts": "",
		name + "/.env":               `DATABASE_URL=`,
		name + "/.gitignore": `node_modules
		env
		`,
		name + "/src/config/db.ts": `import { PrismaClient } from '@prisma/client';
			const prisma = new PrismaClient();
	export default prisma;`,

		//name + "/prisma/schema.prisma":				 ""
	}
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("error while creating dir", err)

		} else {
			fmt.Println("dir created successfully")
		}
	}
	for filepath, content := range files {
		err := CreateFileWithContent(filepath, content)
		if err != nil {
			fmt.Println("error while creating file with content", err)

		}
	}
	cmd := exec.Command("tsc", "--init")
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {

		return err
	}
	//tsConfig file
	tsConfigPath := name + "/tsconfig.json"
	overRideFile(tsConfigPath)
	return nil
}
func overRideFile(filepath string) {
	tsConfig := `{
		"compilerOptions": {
		  "target": "ES6",
		  "module": "commonjs",
		  "outDir": "./dist",
		  "rootDir": "./src",
		  "strict": true,
		  "esModuleInterop": true
		},
		"include": ["src/**/*.ts"],
		"exclude": ["node_modules"]
	  }`
	err := os.WriteFile(filepath, []byte(tsConfig), 0644)
	if err != nil {
		fmt.Printf("Error overriding tsconfig.json: %v\n", err)
		return
	}
	fmt.Println("tsconfig.json overridden successfully!")
}
