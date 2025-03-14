{{- define "sun-panel.conf" -}}
[base]
# Web run port. Default:3002
http_port=3002
# Database driver [mysql/sqlite(Default)]
database_drive=mysql
# Enable static file server. Default:true
enable_static_server=true
# Enable system monitor feature. Default:true
enable_monitor=false
# Used as prefix to generate file url. For example, "/uploads/xxxx.png"
url_prefix=/api/file/s3/

# optional, valid when database_drive=mysql
[mysql]
host={{ .Values.mysql.host}}
port={{ .Values.mysql.port}}
username={{ .Values.mysql.username}}
password={{ .Values.mysql.password}}
db_name={{ .Values.mysql.db_name}}
wait_timeout=100

# Use rclone to store files. Both support local and remote storage
[rclone]
type=s3
provider={{ .Values.rclone.provider}}
access_key_id={{ .Values.rclone.access_key_id}}
secret_access_key={{ .Values.rclone.secret_access_key}}
endpoint={{ .Values.rclone.endpoint}}
region={{ .Values.rclone.region}}
bucket={{ .Values.rclone.bucket}}

[jwt]
# JWT secret key
secret={{ .Values.jwt.secret}}
# JWT token expiration time in hours, default: 72
expire={{ .Values.jwt.expire}}
{{- end }}
