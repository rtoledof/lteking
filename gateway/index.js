const { ApolloServer } = require('apollo-server');
const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require('@apollo/gateway');

class AuthenticatedDataSource extends RemoteGraphQLDataSource {
    willSendRequest({ request, context}) {
        if (context.authorization) {
            request.http.headers.set('Authorization', context.authorization);
        }
    }
}

const config = require('./config.json');

const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
        subgraphs: config.subgraphs
    }),
    buildService({ name, url }) {
        return new AuthenticatedDataSource({ url });
    },
    subscriptions: false
});

const server = new ApolloServer({
    gateway,
    subscriptions: false,
    context: ({ req }) => {
        const authorization = req.headers.authorization || '';
        return { authorization };
    },
});

server.listen().then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`);
});