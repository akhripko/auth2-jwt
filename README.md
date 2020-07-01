# [auth2] multi-provider authentication

In order to load Account Auth Config from back-office-service via gRPC: 
APP_REMOTE_SERVER_TARGET should have an appropriate value (for example service-host:8080) 

By default Config will be loaded from yaml file idpconf.yml
for custom file name use APP_IDP_CONFIG_FILE_NAME

The full list of options can be found here: https://github.com/akhripko/auth2-jwt/blob/master/options/env.go