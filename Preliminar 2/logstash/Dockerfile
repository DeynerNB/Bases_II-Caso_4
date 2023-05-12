FROM logstash:8.7.0

# Instala el plugin de mongodb y cualquier otra dependencia necesaria.
RUN logstash-plugin install logstash-input-mongodb

# Comando para ejecutar Logstash
CMD ["logstash", "-f", "/usr/share/logstash/pipeline/*.conf"]
