const { ApolloServer } = require('apollo-server');
const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require('@apollo/gateway');

class AuthenticatedDataSource extends RemoteGraphQLDataSource {
    willSendRequest({ request, context}) {
        if (context.authorization) {
            request.http.headers.set('Authorization', context.authorization);
        }
    }
}

const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
            { name: 'auth', url: 'http://localhost:3001/query' },
            { name: 'wallet', url: 'http://localhost:3002/query' },
            { name: 'orders', url: 'http://localhost:3003/query' },
        ],
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