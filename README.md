# Stori Challenge Golang

## Descripción del Proyecto

Este proyecto es un sistema para procesar un archivo desde un directorio montado que contiene una lista de transacciones de débito y crédito en una cuenta. El sistema procesa el archivo y envía información resumida al usuario en forma de correo electrónico. Utiliza AWS Lambda, S3 y SES para implementar una solución escalable y serverless. El proyecto está desarrollado en Go y sigue el patrón de diseño de adapters.


## Configuración de AWS Lambda

### Prerrequisitos

- Tener una cuenta en AWS.
- Tener configurado AWS CLI en tu máquina local.
- Tener configurado un bucket S3 para almacenar el archivo de transacciones y el logo.
- Crear roles de IAM con permisos adecuados para cada Lambda function (S3, DynamoDB, SES).

### Configuración de Variables de Entorno

Cada función Lambda necesita ciertas variables de entorno para funcionar correctamente. Asegúrate de configurar las siguientes variables en cada Lambda function:

#### Lambda: read_transactions_s3.zip

1. **Variables de entorno**:
    - `BUCKET_CSV`: Nombre del bucket S3 donde se almacenan los archivos CSV.
    - `DYNAMODB_TABLE`: Nombre de la tabla DynamoDB donde se guardarán las transacciones.

#### Lambda: save_transactions.zip

1. **Variables de entorno**:
    - `DYNAMODB_TABLE`: Nombre de la tabla DynamoDB donde se guardarán las transacciones.

#### Lambda: send_email.zip

1. **Variables de entorno**:
    - `DYNAMODB_TABLE`: Nombre de la tabla DynamoDB donde se guardarán las transacciones.
    - `SES_SENDER`: Dirección de correo electrónico verificada por SES que enviará los correos.
    - `SES_RECIPIENT`: Dirección de correo electrónico del destinatario.
    - `LOGO_URL`: URL del logo que se incluirá en el correo electrónico.

### Despliegue de las Lambdas

1. **Subir los archivos zip a S3** (opcional pero recomendado para despliegue):
    ```sh
    aws s3 cp read_transactions_s3.zip s3://your-bucket-name/read_transactions_s3.zip
    aws s3 cp save_transactions.zip s3://your-bucket-name/save_transactions.zip
    aws s3 cp send_email.zip s3://your-bucket-name/send_email.zip
    ```

2. **Crear las funciones Lambda**:

    - **read_transactions_s3**:
        ```sh
        aws lambda create-function \
            --function-name read_transactions_s3 \
            --runtime go1.x \
            --role arn:aws:iam::your-account-id:role/your-lambda-role \
            --handler bootstrap \
            --code S3Bucket=your-bucket-name,S3Key=read_transactions_s3.zip \
            --environment Variables={S3_BUCKET=your-bucket,DYNAMODB_TABLE=your-dynamodb-table}
        ```

    - **save_transactions**:
        ```sh
        aws lambda create-function \
            --function-name save_transactions \
            --runtime go1.x \
            --role arn:aws:iam::your-account-id:role/your-lambda-role \
            --handler bootstrap \
            --code S3Bucket=your-bucket-name,S3Key=save_transactions.zip \
            --environment Variables={DYNAMODB_TABLE=your-dynamodb-table}
        ```

    - **send_email**:
        ```sh
        aws lambda create-function \
            --function-name send_email \
            --runtime go1.x \
            --role arn:aws:iam::your-account-id:role/your-lambda-role \
            --handler bootstrap \
            --code S3Bucket=your-bucket-name,S3Key=send_email.zip \
            --environment Variables={SES_SENDER=your-sender-email,SES_RECIPIENT=your-recipient-email,LOGO_URL=your-logo-url}
        ```

3. **Configurar triggers**:
    - Para `read_transactions_s3`: Configurar el trigger para que se active al subir un archivo CSV al bucket S3.
    - Para `save_transactions`: Configurar el trigger para que se active con eventos SQS o DynamoDB Streams.
    - Para `send_email`: Configurar el trigger para que se active con eventos SQS o DynamoDB Streams.

## Uso

1. **Sube el archivo CSV al bucket S3**.
2. **La función Lambda `read_transactions_s3` se activa** y procesa el archivo.
3. **Las transacciones se guardan en DynamoDB** usando la función `save_transactions`.
4. **La función `send_email` se activa** y envía un resumen de las transacciones por correo electrónico.

## participante del challenge

Jaime Vallejo 
[Linkedin](https://www.linkedin.com/in/jaime-daniel-vallejo-pe%C3%B1a-37586a170/)

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.
