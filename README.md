**Name: Ananth Upadhya**

**SJSU ID: 015234726**

**Project -1**

**Requirements:**

Create a highly available, highly scalable, cost effective 3 tier web application which would be accessible over public internet through your registered domain name. E.g. [www.mysite.com/project-1](http://www.mysite.com/project-1)

My Website :- [http://fileexplorerpost.xyz/](http://fileexplorerpost.xyz/)

**Steps to Run**
Backend :- export the necessary variables (DynamoDB variables)
```
#!/bin/bash
export DB_AWS_ACCESS_KEY_ID=
export DB_AWS_SECRET_ACCESS_KEY=
export DB_AWS_REGION=
export DYNAMODB_ENDPOINT_URL=
export S3_AWS_DEFAULT_REGION=
export S3_TITAN_BUCKET_NAME=
export S3_AWS_ACCESS_KEY_ID=
export S3_AWS_SECRET_ACCESS_KEY=


go run main.go

```
Frontend :- 

```
npm install
ng serve
```

**Project 1 :-**

a)Criteria:- **User Registration**

Points :- 5 points

User registration form has fields first name, last name, email address, isAdmin (Check and Uncheck), Password. The user should register first to successfully use the app.

![alt text](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image2.png)

b)Criteria :- **Custom Login**

Points :- 5 points

Once registration is complete the user can login to the portal by clicking on login button. The user has to insert the his/her email address and then password and the after successful login the user is taken into the file page.

![alt_text](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image3.png)

c) Criteria :- **File Upload**

Points :- 5 points

Here the user can see the files which he/she has uploaded and can see the details of the file like file name, file description, file created at and file updated at fields.

The users will see the files which belongs to him/her. If the user is an admin user, then the user will be able to see the admin icon clicking which he will be able to see all the files uploaded by all the users.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image4.png)

Users can upload any type of files with the only restriction that files greater than 10MB cannot be uploaded.

d) Criteria :- **File Download**

Points :- 5 points

Users can download these files which they uploaded by clicking on the download icon.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image5.png)

e)Criteria :- **Database Updates**

Points :- 5 Points

I am using Amazon DynamoDB for this project. I use Query, Scan and Insert/Update/Delete Item to perform Read, Insert, Update and Delete items/attributes in the Database

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image6.png)

f)Criteria :- **File Edit**

Points :- 5 Points

Users can edit these files, they can update the file description, file name and the file. Updated values of these fields are shown in the portal.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image7.png)

g)Criteria :- **File Delete**

Points :- 5 Points

Users/Admin can delete the files using the delete button.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image8.png)

**AWS Configurations and Usages**

h)Criteria :- **R53**

Points :- 5 Points

I purchased a domain in Godaddy as the cost of the domain was cheaper around ($1-$2) when compared to AWS domains which costed around ($10-$12). Once the domain was obtained ([http://fileexplorerpost.xyz](http://fileexplorerpost.xyz/)) using amazon R53 created a hosted zone. After creating a hosted zone added a record with the domain name and EC2 public IP. As the domain name was purchased in GoDaddy, I added amazon nameservers to the GoDaddy account.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image9.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image10.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image11.png)

i)Criteria: **ELB (Elastic Load Balancer)**

Points :- 5 Points

I have used Classic Elastic Load Balancer for distributing the traffic across multiple ec2 instances. I defined the load balancer, assigned the security groups, configured the healthz check and registered EC2 instances with load balancer. I have my ec2 instances at different availability zones and have enabled the cross zone load balancing.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image12.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image13.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image14.png)

j)Criteria :- **S3 Bucket and Cloud Front**

Points :- 5 Points

I have created an S3 bucket where the files are getting uploaded. I have used unique keys to store these files, there by increasing performance. I have created another s3 bucket in a different region to have a disaster recovery in place. Apart from this I have done all the necessary configuration which was required for assignment 2. i.eI have used bucket-defined lifecycle strategy. Doing this I can transfer the object to various storage classes such as SIA and Glacier at configured time intervals as the cost can be reduced, I have used SIA to transfer infrequent access data, but it will be instant when required. The archive data that can be used for auditing or other year-end events can be shifted by Glacier. Following the policy I have specified on the buckets as per the current requirement, after 75 days from the date of creation, I will transfer the objects from Standard S3 to SIA and then I will hold it there for a few months until a year is completed from the date of creation of the object. Once a year is completed, I am keeps the data in the glacier for another year that can be used for audit and then we purge the data. (i.e 365 \* 2 = 730 days).

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image15.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image16.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image17.png)

**Cloud Front**

AWS provides a CDN offering called cloud front. Cloud Front can be designed to support the S3 bucket content so that data can be easily served across the globe via the globally deployed edge location. There can be reliability due to various edge locations. The cloud front delivery and its configuration are seen in the following screenshot. Also enabling transfer acceleration helps in faster data transfers from s3 bucket.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image18.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image19.png)

j)Criteria :- **Lambda**

Points :- 5 Points

I have used AWS lambda to compress the files when push event to S3 bucket occurs. This will ensure that the backup is using less of storage space. As the compressed file is only used for backup the time it takes to decompress can be ignored as we don&#39;t do any real time retrieval from this backup bucket.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image20.png)

g)Criteria :- **SNS, Cloudwatch**

Points :- 5 Points

Simple Notification Service is being used to notify the default admin/user whenever lambda function fails to compress the files. An email notification would be triggered to the user. I have created a topic and subscription. So I have configured my email address for the same.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image21.png)

Cloud Watch :

I am using amazon cloud watch for metrics. My application is running on ec2 Instances and the various metrics from these instances are sent to cloudwatch which gives a neat overview on various data like CPU Utilization, Memory Consumpution, Requests coming into across time interval, Which all requests failed and what were the status code of these request/responses. My s3 bucket is also configured with cloud watch which gives me data to storage. I have also integrated my Lambda function with cloud watch which gives me real time statics.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image22.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image23.png)

g)Criteria :- **DR Measures**

Points :- 5 Points

Building the disaster site in different region of a different geographical area allows the best use of the disaster site, so I built two buckets in two different areas, one in US East and one in Asia Pacific. The replication function of AWS S3 will help us in copying the data from the bucket of one region to other region.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image24.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image25.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image26.png)

i)Criteria :- **High Availability Solution (Multi AZ Replication)**

Points :- 5 Points

I have configured my aws resources/services to be highly available most of the time. I have used DynamoDB which by default is Multi AZ Enabled. The S3 bucket which I have primary and backup are in different Availability Zone. Also my EC2 instances are in different AZ&#39;s. And the ELB ensures that the traffic flows across these EC2 instances spread across different Availability Zones.

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image27.png)

j) Criteria :-**Highly Scalable (Autoscale Group)**

Points :- 10 Points

I have enabled AutoScaling for my EC2 instances. I have specified the scale in and scale out rule based on CPU Utilization and memory consumption. This enables a highly scalable solution

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image28.png)

k)Criteria :- **Version Control Github, Codestar, CodeCommit, other**

Points :- 10 Points

I have used Github for maintaining my code. I have also used Amazon Code commit which is also similar to Github. Amazon CodeCommit has features to create branches, raise pull requests, maintain the code commits, maintain versions etc. Amazon Codestar helps in building a pipeline for the projects which comes very handy for easier development. I have written the backend completely in GoLang. I have used Amazon GoLang sdk to communicate with Amazon Web Services. I have followed the best practices while writing my code in GoLang. I have taken care of the error scenarios, corner cases etc.

Frontend is written in Angular, I have created a simple front-end to talk to my backend server.

Github URL :- [https://github.com/ANANTHUPADHYA/cloud](https://github.com/ANANTHUPADHYA/cloud)

CodeCommit

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image29.png)

l)Criteria :- **Admin Panel**

Points :- 10 Points

Admin panel lists all the files and the details from all the users. Admin can view these files from all the users who are using the application and admin also has the capability to delete the files. A small button is displayed at the top of the screen, clicking on which an admin panel will open. The

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image30.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image31.png)

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image32.png)

j)Criteria :- **UI, Documentation, Video, AWS Resource Config**

Points :- 10 Points

UI has been developed using angular and the actual UI can be seen once we go the domain. ([http://fileexplorerpost.xyz/](http://fileexplorerpost.xyz/))

Documentation of the same is written here.

AWS Resource Config

![](https://github.com/ANANTHUPADHYA/cloud/blob/master/screenshots/image33.png)

Thank you
