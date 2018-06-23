
Based on https://www.vultr.com/docs/how-to-install-rabbitmq-on-ubuntu-16-04-47

How To Install RabbitMQ on Ubuntu 16.04
Published on: Fri, Mar 16, 2018 at 11:28 am EST
Linux Guides Server Apps Ubuntu
This article is a port of "How to Install RabbitMQ on CentOS 7" for Ubuntu 16.04.
RabbitMQ is a widely used open-source message broker written in the Erlang programming language. As a message-oriented middleware, RabbitMQ can be used to implement the Advanced Message Queuing Protocol (AMQP) on all modern operating systems.

This article explains how to install RabbitMQ on a Vultr Ubuntu 16.04 server instance.

== Prerequisites

Before getting started, you need to:

* Deploy a brand new Vultr Ubuntu 16.04 server instance.

* Log into the server as a non-root sudo user.


==  Step 1: Update the system
Use the following commands to update your Ubuntu 16.04 system to the latest stable status:

```
sudo apt-get update
sudo apt-get upgrade
```

==  Step 2: Install Erlang
Since RabbitMQ is written in Erlang, you need to install Erlang before you can use RabbitMQ:

```
cd ~
wget http://packages.erlang-solutions.com/site/esl/esl-erlang/FLAVOUR_1_general/esl-erlang_20.1-1~ubuntu~xenial_amd64.deb
sudo dpkg -i esl-erlang_20.1-1\~ubuntu\~xenial_amd64.deb
Verify your installation of Erlang:
```

```
erl
```

You will be brought into the Erlang shell which resembles:

```
Erlang/OTP 20 [erts-9.1] [source] [64-bit] [smp:8:8] [ds:8:8:10] [async-threads:10] [hipe] [kernel-poll:false]

Eshell V9.1  (abort with ^G)
Press Ctrl+C twice to quit the Erlang shell.
```

==  Step 3: Install RabbitMQ

Add the Apt repository to your Apt source list directory (/etc/apt/sources.list.d):

```
echo "deb https://dl.bintray.com/rabbitmq/debian xenial main" | sudo tee /etc/apt/sources.list.d/bintray.rabbitmq.list
```

Next add our public key to your trusted key list using apt-key:

```
wget -O- https://www.rabbitmq.com/rabbitmq-release-signing-key.asc | sudo apt-key add -
```


Run the following command to update the package list:

```
sudo apt-get update
```

Install the rabbitmq-server package:

```
sudo apt-get install rabbitmq-server
```


== Step 4: Start the Server

```
sudo systemctl start rabbitmq-server.service
sudo systemctl enable rabbitmq-server.service
```

You can check the status of RabbitMQ:

```
sudo rabbitmqctl status
```

By default RabbitMQ creates a user named "guest" with password "guest”. You can also create your own administrator account on RabbitMQ server using following commands. Change password to your own password.

```
sudo rabbitmqctl add_user admin password 
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

== Step 5: Enable and use the RabbitMQ management console

Enable the RabbitMQ management console so that you can monitor the RabbitMQ server processes from a web browser:

```
sudo rabbitmq-plugins enable rabbitmq_management
sudo chown -R rabbitmq:rabbitmq /var/lib/rabbitmq/
```

Next, you need to setup an administrator user account for accessing the RabbitMQ server management console. In the following commands, "mqadmin" is the administrator's username, "mqadminpassword" is the password. Remember to replace them with your own.

```
sudo rabbitmqctl add_user mqadmin mqadminpassword
sudo rabbitmqctl set_user_tags mqadmin administrator
sudo rabbitmqctl set_permissions -p / mqadmin ".*" ".*" ".*"
```
Now, visit the following URL:

http://[your-vultr-server-IP]:15672/
Log in with the credentials you had specified earlier. You will be greeted with the RabbitMQ remote management console, where you can learn more about RabbitMQ.