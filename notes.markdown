 kind delete cluster --name trailbox                                 
kind create cluster --name trailbox
kubectl get nodes           
 kubectl apply -f k8s/postgres.yaml       
kubectl get pods                                                    
kubectl get svc
kubectl apply -f k8s/users-http-deployment.yaml                     
kubectl apply -f k8s/users-http-service.yaml                        
kubectl get pods           
kubectl exec -it deployment/postgres -- psql -U trailbox -d trailbox
    DROP TABLE IF EXISTS users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT,
  email TEXT
);

INSERT INTO users (name, email)
VALUES 
('Samuel Martínez', 'sam@example.com'),
('Eduardo Saenz', 'chingon@example.com');

kubectl port-forward service/trailbox-users-http-service 8080:80  