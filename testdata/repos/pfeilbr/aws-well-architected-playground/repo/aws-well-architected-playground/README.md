# aws-well-architected-playground

deep dive on all things AWS Well-Architected

## Key Points

* consistent pre-launch review process against AWS best practices
* helps you understand the pros and cons of decisions you make while building systems
* review process is a conversation and not an audit.  working together to improve. practical advice.
* goal is not to have "perfect" architecture from the start. identify areas for improvement and choose a couple that delivery the most value
* AWS does not provide prescriptive guidance on how to perform the review.  WA tool is the closest.
* concepts: Pillars -> Design Principles -> Questions
* enables: learn -> measure -> improve iterative cycle
* input: answer questions, output: improvement plan (PDF reports)
* learning / education - can be used as standalone tool solely for learning what the best practices are
* milestone - record the state of a workload for given point in time. e.g. original design, design review, v1, v2

## Use Cases

* learning best practices for the cloud
* technology governance
* portfolio management - inventory of workloads, historical decisions made, risks, highlights where to invest

## Well-Architected Framework

> The AWS Well-Architected Framework helps you to design and operate a reliable, secure, efficient, and cost-efficient systems on AWS. It also helps you constantly measure your architecture against best practices and provides you an opportunity to improve your architecture.

### 5 Pillars

* Operational Excellence
* Security
* Reliability
* Performance Efficiency
* Cost Optimization

### Review Process

> The review process describes in high-level terms, how the assessment of the principles should be done. For AWS, this should be a lightweight process, which is taking rather hours, instead of days and it should be repeated multiple times across the architecture lifecycle. AWS states that it is important to have a conversation (not an audit) and a “blame-free approach” during the review and it is also important to involve the right people. The results of the conversations should be a list of issues that can then be prioritized based on the business context and that can be formulated into a set of actions that help to improve the overall customer experience of the architecture.

## Well-Architected Tool

AWS Console Tool that steps a user through the Well-Architected Review Process

### Feature Request

One area where there is a gap for an enterprises are all the company specific policies, standards, and best practices that are additive and need to be addressed on top of AWS.  These types of questions and guidance would need to happen outside of WA Tool.

A feature to define custom lenses - a customer defined lens.  This way the single WA Tool could be the method for review facilitation, improvement reporting and maintaining history.

## Key Visuals

<img src="https://www.evernote.com/l/AAGT8fHSJRtGV4338tvJ0FHMvbb_vfTh7qkB/image.png" alt="WA Tool Features" width="50%" />

<img src="https://www.evernote.com/l/AAE0ZIXQoM5KYK9TUKsBpNZ_dSQ_x27Ti_0B/image.png" alt="General Design Principles" width="50%" />

<img src="https://www.evernote.com/l/AAGxCfQWxsZMDZLQC9dqcm-EJRgkM3jsNlwB/image.png" alt="Intent of WA Review" width="50%" />

<img src="https://www.evernote.com/l/AAF2yN9IdGZOS6_CoSnpzjc0nWOru88e4a0B/image.png" alt="Review Benefits" width="50%" />

<img src="https://www.evernote.com/l/AAEqVMFtRnJBMp-iX7F3Y8eLLIu3kV6ruvQB/image.png" alt="" width="50%" />

## Resources

* [AWS Well-Architected](https://aws.amazon.com/architecture/well-architected)
* [Documentation | AWS Well-Architected Framework](https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html)
* [The Review Process - AWS Well-Architected Framework](https://wa.aws.amazon.com/wat.thereviewprocess.wa-review.en.html)
* [AWS Well Architected framework: A Complete Checklist](https://www.rapyder.com/blogs/aws-well-architected-framework-checklist/)
* [AWS Well-Architected Framework Cheatsheet | Cloud Noon](https://cloudnoon.com/blog/aws/aws-well-architected-framework-cheatsheet/)