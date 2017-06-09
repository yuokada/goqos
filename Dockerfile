FROM centos:centos7

#RUN yum install -y http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-8.noarch.rpm
RUN yum groupinstall -y -q "Development tools"
RUN yum install -y -q python-pip curl pwgen git wget
RUN yum update -y && \
    yum install -y rpmdevtools python2-devel python-sphinx libyaml-devel \
                   gcc make python-setuptools \
                   vim tree golang && \
    yum clean all
RUN wget -q -O go.tar.gz https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
RUN tar zxf go.tar.gz -C /usr/local/

ENV PATH /usr/local/go/bin:$PATH
RUN mkdir -p /rpmbuild && mkdir -p /root/rpmbuild


COPY ./                 /rpmbuild/
COPY ./.rpmmacros       /root/

RUN chown root:root -R /rpmbuild
WORKDIR /rpmbuild
RUN /bin/bash -x  ./buildrpm.sh
RUN rpm -qpl rpmbuild/RPMS/x86_64/goqos-0.1-0.el7.centos.x86_64.rpm

#RUN /usr/bin/easy_install-2.7 pip && pip2.7 install pypi2rpm
#RUN  for f in `find SPECS -name "*.spec"` ; do rpmbuild -ba ${f}; done
CMD /bin/bash
