# Establecer la imagen base
FROM node:latest

# Crear directorio de la aplicación
WORKDIR /usr/src/app

# Instalar las dependencias de la aplicación
# Una wildcard se utiliza para asegurar que tanto package.json como package-lock.json sean copiados
# donde están disponibles (npm@5+)
COPY package*.json ./

RUN npm install
# Si estás construyendo tu código para producción
# RUN npm ci --only=production

# Empaquetar el código de la aplicación
COPY . .

# Exponer el puerto en el que se ejecutará el servidor
EXPOSE 4000

# Comando para iniciar la aplicación
CMD [ "node", "index.js" ]