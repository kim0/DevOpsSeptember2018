# Copyright 2018 Ahmed Kamal

FROM scratch
COPY ./virus /virus
CMD [ "/virus" ]