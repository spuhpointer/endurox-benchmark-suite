Install based on https://www.virtual-server.org/blog/how-to-install-activemq-on-ubuntu-16-04/


In this tutorial, we’ll show you, How to install ActiveMQ on Ubuntu 16.04. ctiveMQ is an open-source, message broker software, which can serve as a central point for communication between distributed processes. In other words, it’s used to reliably communicate between two distributed processes or applications/services that need to communicate with each other, but their languages are incompatible between each other. Since its creation, ActiveMQ has turned into a strong alternative to the commercial alternatives like EMS/TIBCO or WebSphereMQ and it is actively used in production by some of the top companies ranging from financial services to retail. Installing ActiveMQ on Ubuntu 16.04, if pretty easy task, just follow the steps bellow and you should have it installed in less then 10 minutes. There is very little information on how ActiveMQ is installed. In reality, the installation is pretty straight forward, but and the official documentation is very vague and this leads to frustration in an inexperienced person. In this blog post, we will explain how ActiveMQ can be installed on a Ubuntu 16 VPS.


First things update/upgrade your Virtual Server with the latest software:

```
# apt-get update && apt-get upgrade
```

# Install the default Java JDK

```
# apt-get install default-jdk
```

# Check if the java has correctly been installed:

```
# java -version
openjdk version "1.8.0_111"
OpenJDK Runtime Environment (build 1.8.0_111-8u111-b14-2ubuntu0.16.04.2-b14)
OpenJDK 64-Bit Server VM (build 25.111-b14, mixed mode)
We are now ready to install ActiveMQ. We will be installing into the /opt directory by downloading the tar archive from their official site. At the moment of installing 5.14.3 is the current version:
```

```
# cd /opt
# wget https://archive.apache.org/dist/activemq/5.15.4/apache-activemq-5.15.4-bin.tar.gz

```

# Unpack the archive and create a simlink to the current installation for easier access later.

```
# tar zxf apache-activemq-5.15.4-bin.tar.gz
# ln -s /opt/apache-activemq-5.15.4 activemq
# rm apache-activemq-5.15.4-bin.tar.gz
```

At this moment the installation is complete. ActiveMQ configuration files are located in its default directory in the conf folder. In our case, they are located into /opt/activemq/conf/. The default configuration files work well for testing the application. The file of our interest is the /opt/activemq/conf/activemq.xml file because in here are defined the transport connectors i.e. the protocols which we like to be enabled or disabled for our ActiveMQ implementation.

Feel free to check the <transportConnectors> section and comment out any unwanted protocols. For testing purposes, you can leave the default as it is.

```
<transportConnectors>
 <!-- DOS protection, limit concurrent connections to 1000 and frame size to 100MB -->
 <transportConnector name="openwire" uri="tcp://0.0.0.0:61616?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
 <transportConnector name="amqp" uri="amqp://0.0.0.0:5672?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
 <transportConnector name="stomp" uri="stomp://0.0.0.0:61613?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
 <transportConnector name="mqtt" uri="mqtt://0.0.0.0:1883?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
 <transportConnector name="ws" uri="ws://0.0.0.0:61614?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
</transportConnectors>
```

Start the ActiveMQ server whenever you feel ready:

```
# /opt/activemq/bin/activemq start
```

# Check that everything is good using the netstat command:

```
# netstat -tulnp | grep java
tcp 0 0 0.0.0.0:1883 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:8161 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:5672 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:42699 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:61613 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:61614 0.0.0.0:* LISTEN 25074/java
tcp 0 0 0.0.0.0:61616 0.0.0.0:* LISTEN 25074/java
```

ActiveMQ’s web administration interface listens on port :8161. It can be accessed using the bellow URL for the front-end and the back-end respectively where WW.XX.YY.ZZ is your VPS IP Address. The default username and password are both admin

http://WW.XX.YY.ZZ:8161
http://WW.XX.YY.ZZ:8161/admin
And that’s it, you now have functioning ActiveMQ service

PS. If you liked this post, on howw to install ActiveMQ on Ubuntu 16.04, please share it with your friends on the social networks using the buttons on the left or simply leave a reply below. Thanks.


# Trouble shooter

If you have RabbitMQ running you may get following error on startup:

```
2018-06-26 22:49:26,449 | ERROR | Failed to start Apache ActiveMQ (localhost, ID:jimbopc-46075-1530042566259-0:1) | org.apache.activemq.broker.BrokerService | main
java.io.IOException: Transport Connector could not be registered in JMX: java.io.IOException: Failed to bind to server socket: amqp://0.0.0.0:5672?maximumConnections=1000&wireFormat.maxFrameSize=104857600 due to: java.net.BindException: Address already in use (Bind failed)
        at org.apache.activemq.util.IOExceptionSupport.create(IOExceptionSupport.java:28)[activemq-client-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.registerConnectorMBean(BrokerService.java:2264)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.startTransportConnector(BrokerService.java:2744)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.startAllConnectors(BrokerService.java:2640)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.doStartBroker(BrokerService.java:771)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.startBroker(BrokerService.java:733)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.broker.BrokerService.start(BrokerService.java:636)[activemq-broker-5.15.4.jar:5.15.4]
        at org.apache.activemq.xbean.XBeanBrokerService.afterPropertiesSet(XBeanBrokerService.java:73)[activemq-spring-5.15.4.jar:5.15.4]
        at sun.reflect.NativeMethodAccessorImpl.invoke0(Native Method)[:1.8.0_171]
```

Thus shutdown of RabbitMQ is needed before trying to start ActiveMQ again:

```
# sudo systemctl stop rabbitmq-server.service
# /opt/activemq/bin/activemq start
```

