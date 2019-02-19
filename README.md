# 1. K8s setting
## Install on Windows.
I followed http://dockone.io/article/8136 to install on my computer. But I meet a lot of problem such as hyper-v issue and etc. So I give up to install K8s on Windows.
## Use GKE (Google Kubernetes Engine)
I follwed the tutorial and use free-trial account to lauch a GKE on Google Cloud Platform. I create a g1-small instance with 1 node in asia-east1-a at first.
# 2. OpenFaas setting
Follow https://github.com/stefanprodan/openfaas-gke to install OpenFaas. But I could not launch OpenFaas at the final step. After check log and status of OpenFaas pod, it shows insufficient memory. It is caused by reserving g1-smal instance which is too small for OpenFaas. After reserving n1-instance-2, OpenFaas could be launched.
# 3. MySQL setting
Follow https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/ to install MySQL server on GKE. MySQL server could be connected via kubectrl exec after install, but could not connected from other site. After some study, I fix the service type of MySQL from Nodeport to LoadBalancer and related setting. MySQL server finnaly could be access from other site.
# 4. GOlang function build up
I followd https://blog.alexellis.io/serverless-golang-with-openfaas/ to add my first Golang function on OpenFaas.
  1. faas-cli new --lang go order //Add new function from template.
  2. Fix gateway to my environment.
  3. Fix code.
  4. faas-cli build -f order.yml //Build function
  5. faas-cli push -f order.yml //Push to docker hub because I am using remote cluster.
  6. faas-cli deploy -f order.yml //Deploy function to server.
  
Currently I could not find any way to pass url path parameter to function. I would pass the parameters from body in 1st step. After competing those functions, I would keep find how to fix URL path.
# 5. Connect to MySQL with Golang
I followed https://blog.alexellis.io/serverless-golang-with-openfaas/ to add mysql liberaries. But build fails due to could not found mysql liberaries. After some testing, it caused by GOPATH setting. Default setting contains some not used path in GKE, it should be export to the function folder.
# 6. Golang function impementation
After study json and map usage in Golang, functions could be run in GKE now.
You could access server by http://34.80.10.134/function/order and http://34.80.10.134/function/inventory.
Here is some snapshot.
* POST order:
![image](https://github.com/andychen120/PentiumQuiz/blob/master/pictures/post_order.png)
* DELETE order
![image](https://github.com/andychen120/PentiumQuiz/blob/master/pictures/delete_order.png)
* GET order
![image](https://github.com/andychen120/PentiumQuiz/blob/master/pictures/get_order.png)
* inventory
![image](https://github.com/andychen120/PentiumQuiz/blob/master/pictures/inventory.png)
# 7. How to use these functions
  1. git clone from repository.
  2. use faas-cli build the yml file
  3. use faas-cli push yml file when using remote cluster
  4. use faas-cli deploy yml file
  5. Try on OpenFaas UI or curl.
# 8. Uncompleted items
* URL path parameter
* Parameter check
* Error handling
