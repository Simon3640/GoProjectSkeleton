package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// Obtener el directorio base (donde están las funciones)
	// Por defecto, buscamos en tmp/ dentro del directorio functions/
	baseDir := "."
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	} else {
		// Si no se especifica, intentar detectar desde el directorio actual
		wd, err := os.Getwd()
		if err == nil {
			// Si estamos en cmd/deploy, subir dos niveles
			if strings.HasSuffix(wd, "cmd/deploy") {
				baseDir = filepath.Join(wd, "../..")
			} else if strings.HasSuffix(wd, "functions") {
				baseDir = wd
			}
		}
	}

	// Siempre verificar si existe tmp/ y usarlo como baseDir si contiene funciones
	// Esto funciona incluso si se pasa "." como argumento
	tmpDir := filepath.Join(baseDir, "tmp")
	if _, err := os.Stat(tmpDir); err == nil {
		// Verificar que tmp/ contiene funciones (buscar al menos un directorio conocido)
		testPath := filepath.Join(tmpDir, "status/health_check")
		if _, err := os.Stat(testPath); err == nil {
			baseDir = tmpDir
		}
	}

	// Obtener Function App name desde variables de entorno o argumentos
	functionAppName := os.Getenv("FUNCTION_APP_NAME")
	if functionAppName == "" && len(os.Args) > 2 {
		functionAppName = os.Args[2]
	}
	if functionAppName == "" {
		fmt.Fprintf(os.Stderr, "Error: FUNCTION_APP_NAME no está configurado\n")
		fmt.Fprintf(os.Stderr, "Uso: go run deploy.go [baseDir] [FUNCTION_APP_NAME]\n")
		fmt.Fprintf(os.Stderr, "O configura la variable de entorno FUNCTION_APP_NAME\n")
		os.Exit(1)
	}

	// Obtener resource group (opcional)
	resourceGroup := os.Getenv("RESOURCE_GROUP")
	if resourceGroup == "" && len(os.Args) > 3 {
		resourceGroup = os.Args[3]
	}

	fmt.Printf("🚀 Iniciando despliegue de funciones a: %s\n", functionAppName)
	if resourceGroup != "" {
		fmt.Printf("📦 Resource Group: %s\n", resourceGroup)
	}
	fmt.Printf("📁 Directorio base: %s\n", baseDir)
	fmt.Println()

	// Obtener todas las funciones
	allFunctions := GetAll()

	// Obtener el directorio de trabajo actual
	initialWd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error obteniendo directorio de trabajo: %v\n", err)
		os.Exit(1)
	}

	// Convertir baseDir a ruta absoluta
	if !filepath.IsAbs(baseDir) {
		baseDir = filepath.Join(initialWd, baseDir)
	}
	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error obteniendo ruta absoluta: %v\n", err)
		os.Exit(1)
	}

	// Mutex para sincronizar la salida
	var outputMutex sync.Mutex

	// Paso 1: Compilar todas las funciones en paralelo
	fmt.Printf("📦 Compilando todas las funciones...\n")
	fmt.Println()

	type compileResult struct {
		path  string
		err   error
		index int
	}
	compileResults := make(chan compileResult, len(allFunctions))
	var compileWg sync.WaitGroup

	for i, fn := range allFunctions {
		compileWg.Add(1)
		go func(index int, fn Config) {
			defer compileWg.Done()

			funcDir := filepath.Join(baseDir, fn.Path)
			funcDir, _ = filepath.Abs(funcDir)

			// Verificar que el directorio existe
			if _, err := os.Stat(funcDir); os.IsNotExist(err) {
				outputMutex.Lock()
				fmt.Printf("⚠️  [%d/%d] Directorio no encontrado: %s\n", index+1, len(allFunctions), fn.Path)
				outputMutex.Unlock()
				compileResults <- compileResult{path: fn.Path, err: fmt.Errorf("directorio no encontrado"), index: index}
				return
			}

			outputMutex.Lock()
			fmt.Printf("[%d/%d] 🔨 Compilando: %s\n", index+1, len(allFunctions), fn.Path)
			outputMutex.Unlock()

			err := compileFunction(funcDir, fn.Path, &outputMutex)
			compileResults <- compileResult{path: fn.Path, err: err, index: index}
		}(i, fn)
	}

	go func() {
		compileWg.Wait()
		close(compileResults)
	}()

	// Recopilar resultados de compilación
	type failedCompilation struct {
		path string
		err  error
	}
	compiledFunctions := make([]Config, 0)
	failedCompilations := make([]failedCompilation, 0)

	for res := range compileResults {
		if res.err != nil {
			failedCompilations = append(failedCompilations, failedCompilation{
				path: res.path,
				err:  res.err,
			})
		} else {
			compiledFunctions = append(compiledFunctions, allFunctions[res.index])
		}
	}

	if len(failedCompilations) > 0 {
		fmt.Printf("\n❌ Error compilando %d funciones:\n", len(failedCompilations))
		for _, failed := range failedCompilations {
			fmt.Printf("   - %s: %v\n", failed.path, failed.err)
		}
		os.Exit(1)
	}

	if len(compiledFunctions) == 0 {
		fmt.Printf("❌ No hay funciones compiladas para desplegar\n")
		os.Exit(1)
	}

	fmt.Printf("\n✅ Todas las funciones compiladas exitosamente (%d funciones)\n", len(compiledFunctions))
	fmt.Println()

	// Paso 2: Crear un paquete único con todas las funciones y desplegar
	fmt.Printf("📦 Creando paquete de despliegue único...\n")
	err = deployAllFunctions(baseDir, compiledFunctions, functionAppName, resourceGroup, &outputMutex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error desplegando: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Todas las funciones desplegadas exitosamente (%d funciones)\n", len(compiledFunctions))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Resumen del despliegue:\n")
	fmt.Printf("   ✅ Funciones desplegadas: %d\n", len(compiledFunctions))
	fmt.Printf("   📦 Total: %d\n", len(allFunctions))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func compileFunction(funcDir, funcPath string, outputMutex *sync.Mutex) error {
	// Convertir a ruta absoluta
	funcDir, err := filepath.Abs(funcDir)
	if err != nil {
		return fmt.Errorf("error obteniendo ruta absoluta: %w", err)
	}

	// 1. Verificar que main.go existe
	mainGoPath := filepath.Join(funcDir, "main.go")
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return fmt.Errorf("main.go no encontrado en %s", funcDir)
	}

	// 2. Ejecutar go mod tidy (usando Dir para especificar el directorio)
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = funcDir
	var tidyOut strings.Builder
	cmd.Stdout = &tidyOut
	cmd.Stderr = &tidyOut
	if err := cmd.Run(); err != nil {
		outputMutex.Lock()
		fmt.Printf("  [%s] ⚠️  Warning en go mod tidy: %s\n", funcPath, tidyOut.String())
		outputMutex.Unlock()
		// Continuar aunque haya warnings
	}

	// Descargar dependencias para generar go.sum
	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = funcDir
	cmd.Stdout = &tidyOut
	cmd.Stderr = &tidyOut
	cmd.Run() // Ignorar errores, go build también descargará

	// 3. Compilar el binario (usando rutas absolutas)
	binDir := filepath.Join(funcDir, "bin")
	os.MkdirAll(binDir, 0755)
	binPath := filepath.Join(binDir, "main")

	cmd = exec.Command("go", "build", "-o", binPath, mainGoPath)
	cmd.Dir = funcDir
	var buildOut strings.Builder
	cmd.Stdout = &buildOut
	cmd.Stderr = &buildOut
	if err := cmd.Run(); err != nil {
		// Mostrar el error de compilación
		outputMutex.Lock()
		fmt.Printf("  [%s] ❌ Error de compilación:\n%s\n", funcPath, buildOut.String())
		outputMutex.Unlock()
		return fmt.Errorf("error compilando: %w", err)
	}

	// Verificar que el binario se creó
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		outputMutex.Lock()
		fmt.Printf("  [%s] ❌ Binario no encontrado en %s\n", funcPath, binPath)
		fmt.Printf("  [%s] Output de compilación: %s\n", funcPath, buildOut.String())
		outputMutex.Unlock()
		return fmt.Errorf("binario no se generó correctamente en %s", binPath)
	}

	return nil
}

func deployAllFunctions(baseDir string, functions []Config, functionAppName, resourceGroup string, outputMutex *sync.Mutex) error {
	// Crear un directorio temporal para el paquete de despliegue
	tempDir, err := os.MkdirTemp("", "azure-functions-deploy-*")
	if err != nil {
		return fmt.Errorf("error creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outputMutex.Lock()
	fmt.Printf("📦 Creando paquete de despliegue en: %s\n", tempDir)
	outputMutex.Unlock()

	// Copiar host.json (usar el primero como referencia, todos deberían ser iguales)
	firstFuncDir := filepath.Join(baseDir, functions[0].Path)
	hostJsonPath := filepath.Join(firstFuncDir, "host.json")
	var cmd *exec.Cmd
	if _, err := os.Stat(hostJsonPath); err == nil {
		cmd = exec.Command("cp", hostJsonPath, tempDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error copiando host.json: %w", err)
		}
	}

	// Copiar cada función: su carpeta con function.json
	// NOTA: En Azure Functions Custom Handler, todas las funciones comparten el mismo binario
	// pero cada una tiene su propia carpeta con function.json
	for _, fn := range functions {
		funcDir := filepath.Join(baseDir, fn.Path)
		functionName := filepath.Base(fn.Path)

		// Crear la carpeta de la función en el paquete
		funcPackageDir := filepath.Join(tempDir, functionName)
		os.MkdirAll(funcPackageDir, 0755)

		// Copiar function.json
		functionJsonPath := filepath.Join(funcDir, functionName, "function.json")
		if _, err := os.Stat(functionJsonPath); err == nil {
			cmd = exec.Command("cp", functionJsonPath, funcPackageDir)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error copiando function.json para %s: %w", fn.Path, err)
			}
		} else {
			return fmt.Errorf("function.json no encontrado para %s en %s", fn.Path, functionJsonPath)
		}
	}

	// Copiar un solo binario compartido (usar el primero como referencia)
	// En un Custom Handler real, todas las funciones usarían el mismo ejecutable
	binPath := filepath.Join(firstFuncDir, "bin", "main")
	if _, err := os.Stat(binPath); err == nil {
		binPackageDir := filepath.Join(tempDir, "bin")
		os.MkdirAll(binPackageDir, 0755)
		cmd = exec.Command("cp", binPath, filepath.Join(binPackageDir, "main"))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error copiando binario: %w", err)
		}
	} else {
		return fmt.Errorf("binario no encontrado en %s", binPath)
	}

	// Crear zip del paquete completo
	zipName := "functions-package.zip"
	zipPath := filepath.Join(tempDir, "..", zipName)
	defer os.Remove(zipPath)

	outputMutex.Lock()
	fmt.Printf("📦 Creando archivo zip...\n")
	outputMutex.Unlock()

	// Cambiar al directorio temporal para crear el zip
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	// Crear zip con todos los archivos
	cmd = exec.Command("zip", "-r", zipPath, ".")
	var zipOut strings.Builder
	cmd.Stdout = &zipOut
	cmd.Stderr = &zipOut
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creando zip: %w", err)
	}

	// Desplegar usando az functionapp deployment source config-zip
	outputMutex.Lock()
	fmt.Printf("☁️  Desplegando paquete completo a Azure Functions...\n")
	outputMutex.Unlock()

	azArgs := []string{"functionapp", "deployment", "source", "config-zip",
		"--name", functionAppName,
		"--src", zipPath,
	}

	if resourceGroup != "" {
		azArgs = append(azArgs, "--resource-group", resourceGroup)
	}

	cmd = exec.Command("az", azArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error desplegando zip: %w", err)
	}

	return nil
}

func buildAndDeployFunction(funcDir, funcPath, functionAppName, resourceGroup string, outputMutex *sync.Mutex) error {
	// Cambiar al directorio de la función
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio actual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(funcDir); err != nil {
		return fmt.Errorf("error cambiando al directorio %s: %w", funcDir, err)
	}

	// 1. Verificar que main.go existe
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		return fmt.Errorf("main.go no encontrado en %s", funcDir)
	}

	// 2. Ejecutar go mod tidy
	outputMutex.Lock()
	fmt.Printf("  [%s] 📥 Ejecutando go mod tidy...\n", funcPath)
	outputMutex.Unlock()

	cmd := exec.Command("go", "mod", "tidy")
	// Capturar output para evitar mezclas
	var tidyOut strings.Builder
	cmd.Stdout = &tidyOut
	cmd.Stderr = &tidyOut
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando go mod tidy: %w", err)
	}

	// 3. Compilar el binario
	outputMutex.Lock()
	fmt.Printf("  [%s] 🔨 Compilando binario...\n", funcPath)
	outputMutex.Unlock()

	os.MkdirAll("bin", 0755)
	cmd = exec.Command("go", "build", "-o", "bin/main", "main.go")
	var buildOut strings.Builder
	cmd.Stdout = &buildOut
	cmd.Stderr = &buildOut
	if err := cmd.Run(); err != nil {
		outputMutex.Lock()
		fmt.Printf("  [%s] ❌ Error compilando: %v\n", funcPath, err)
		outputMutex.Unlock()
		return fmt.Errorf("error compilando: %w", err)
	}

	// Verificar que el binario se creó
	if _, err := os.Stat("bin/main"); os.IsNotExist(err) {
		return fmt.Errorf("binario no se generó correctamente")
	}

	outputMutex.Lock()
	fmt.Printf("  [%s] ✅ Binario compilado exitosamente\n", funcPath)
	outputMutex.Unlock()

	// 4. Desplegar usando func azure functionapp publish
	outputMutex.Lock()
	fmt.Printf("  [%s] ☁️  Desplegando a Azure Functions...\n", funcPath)
	outputMutex.Unlock()

	// Verificar si func está disponible
	if _, err := exec.LookPath("func"); err != nil {
		outputMutex.Lock()
		fmt.Printf("  [%s] ⚠️  'func' no encontrado, usando zip deploy...\n", funcPath)
		outputMutex.Unlock()
		return deployWithZip(funcDir, funcPath, functionAppName, resourceGroup, outputMutex)
	}

	// Construir comando func azure functionapp publish
	args := []string{"azure", "functionapp", "publish", functionAppName, "--force"}

	// Si hay resource group, agregarlo
	if resourceGroup != "" {
		args = append(args, "--resource-group", resourceGroup)
	}

	// No hacer build porque ya compilamos
	args = append(args, "--no-build")

	cmd = exec.Command("func", args...)
	var funcOut strings.Builder
	cmd.Stdout = &funcOut
	cmd.Stderr = &funcOut
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		// Si falla, intentar con zip deploy como respaldo
		outputMutex.Lock()
		fmt.Printf("  [%s] ⚠️  Error con func, intentando zip deploy...\n", funcPath)
		outputMutex.Unlock()
		return deployWithZip(funcDir, funcPath, functionAppName, resourceGroup, outputMutex)
	}

	return nil
}

func deployWithZip(funcDir, funcPath, functionAppName, resourceGroup string, outputMutex *sync.Mutex) error {
	// Crear un zip con los archivos necesarios
	zipName := fmt.Sprintf("%s.zip", strings.ReplaceAll(funcPath, "/", "_"))
	zipPath := filepath.Join(funcDir, zipName)
	defer os.Remove(zipPath)

	outputMutex.Lock()
	fmt.Printf("  [%s] 📦 Creando paquete zip...\n", funcPath)
	outputMutex.Unlock()

	// Verificar que los archivos existen antes de incluirlos
	filesToInclude := []string{}

	// Obtener el nombre de la función (último segmento del path)
	// Ejemplo: "status/health_check" -> "health_check"
	functionName := filepath.Base(funcPath)

	requiredFiles := map[string]string{
		"bin/main":                      "binario",
		"host.json":                     "host.json",
		functionName + "/function.json": "function.json",
	}

	for file, desc := range requiredFiles {
		fullPath := filepath.Join(funcDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("archivo requerido no encontrado: %s (%s)", file, desc)
		}
		filesToInclude = append(filesToInclude, file)
	}

	// Crear comando zip
	zipArgs := []string{"-r", zipName}
	zipArgs = append(zipArgs, filesToInclude...)

	cmd := exec.Command("zip", zipArgs...)
	cmd.Dir = funcDir
	var zipOut strings.Builder
	cmd.Stdout = &zipOut
	cmd.Stderr = &zipOut

	if err := cmd.Run(); err != nil {
		// Si zip no está disponible, intentar con tar
		if strings.Contains(err.Error(), "command not found") {
			return deployWithTar(funcDir, funcPath, functionAppName, resourceGroup, filesToInclude, outputMutex)
		}
		return fmt.Errorf("error creando zip: %w", err)
	}

	// Verificar que el zip se creó
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo zip no se creó correctamente")
	}

	// Desplegar usando az functionapp deployment source config-zip
	outputMutex.Lock()
	fmt.Printf("  [%s] ☁️  Desplegando zip a Azure...\n", funcPath)
	outputMutex.Unlock()

	azArgs := []string{"functionapp", "deployment", "source", "config-zip",
		"--name", functionAppName,
		"--src", zipPath,
	}

	if resourceGroup != "" {
		azArgs = append(azArgs, "--resource-group", resourceGroup)
	}

	cmd = exec.Command("az", azArgs...)
	var azOut strings.Builder
	cmd.Stdout = &azOut
	cmd.Stderr = &azOut

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error desplegando zip: %w", err)
	}

	return nil
}

func deployWithTar(funcDir, funcPath, functionAppName, resourceGroup string, files []string, outputMutex *sync.Mutex) error {
	// Crear un tar.gz con los archivos necesarios
	tarName := fmt.Sprintf("%s.tar.gz", strings.ReplaceAll(funcPath, "/", "_"))
	tarPath := filepath.Join(funcDir, tarName)
	defer os.Remove(tarPath)

	outputMutex.Lock()
	fmt.Printf("  [%s] 📦 Creando paquete tar.gz...\n", funcPath)
	outputMutex.Unlock()

	// Crear comando tar
	tarArgs := []string{"-czf", tarName}
	tarArgs = append(tarArgs, files...)

	cmd := exec.Command("tar", tarArgs...)
	cmd.Dir = funcDir
	var tarOut strings.Builder
	cmd.Stdout = &tarOut
	cmd.Stderr = &tarOut

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creando tar.gz: %w", err)
	}

	// Convertir tar.gz a zip o usar directamente
	// Azure Functions acepta zip, así que necesitamos convertir
	outputMutex.Lock()
	fmt.Printf("  [%s] ⚠️  Azure Functions requiere zip. Por favor instala 'zip' o 'func' para continuar.\n", funcPath)
	outputMutex.Unlock()
	return fmt.Errorf("zip no disponible y tar.gz no es compatible con Azure Functions")
}
