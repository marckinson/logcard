"# logcard" 
Logcard Demo
=======

##Application overview##
This application is designed to demonstrate how assets can be modeled on the Blockchain using a Logcard lifecycle. 
In the scenario parts are modeled using Blockchain technology with the following attributes:

| Attribute       | Type                                                                                                  |
| --------------- | ----------------------------------------------------------------------------------------------------- |
| Id           	  | String																								  |
| PN              | String                      		                                                                  |
| SN              | String                                                                                                |
| PType           | String                                                                                                |
| Owner           | String                                                                                                |
| Responsible     | String                                                                                                |
| signature       | String                                 					                                              |
| Logs            | [] Log                                   	                                                          |

The application is designed to allow participants to interact with the parts assets creating, transferring them as their permissions allow. 
The participants included in the application are as follows:

| Participant    | Permissions                                                          |
| -------------- | ---------------------------------------------------------------------|
| Supplier     	 | Create, Read, Claim, Transfer                     		            |
| AH      	     | Create, Read, Claim, Transfer							            |
| MT_USER        | Read, Transfer														|
| Customer       | Read, Claim, Transfer          								        |
| Certifier      | Read											                   		|
| Shipping 		 | Transfer       						                                |


The demonstration allows a view of the ledger that stores all the interactions that the above participants have has with their assets. 
The ledger view shows the regulator every transaction that has occurred showing who tried to to what at what time and to which part. 
The ledger view also allows the user to see transactions that they were involved with as well as showing the interactions 
with the assets they own before they owned them e.g. they can see when it was created.


##Application scenario##
The scenario goes through the lifecycle of a part which has the following stages:

####Stages:####

 1. Part is created by the Supplier/AH. Supplier/AH is the current owner.
 2. Supplier/AH transfers the part to Shipping_Company which becomes the current responsible of the part.
 3. Shipping_Company transfers the part to customer who becomes the current responsible of the part.
 4. Customer Claims Ownership and becomes the current Owner of the part.
 5. For maintenance reasons, customer transfers the part to MT_USER which becomes the current responsible.
 6. When maintenance is done, MT_USER transfers the part to Customer who become the current responsible.
 7. For Audit need, Authorities (FAA, EASA) review all the historic of concerned parts.