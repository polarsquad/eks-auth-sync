FROM gcr.io/distroless/base
COPY eks-auth-sync /usr/local/bin/eks-auth-sync
ENTRYPOINT ["eks-auth-sync"]
