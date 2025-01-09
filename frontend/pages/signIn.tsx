import { useState } from "react";

import { ViewIcon, ViewOffIcon } from "@chakra-ui/icons";
import {
  Button,
  Center,
  Container,
  Input,
  Link,
  Text,
  VStack,
} from "@chakra-ui/react";

export default function Login() {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <Container>
      <Center>
        <VStack mt={8}>
          <Text fontSize="2xl">Login</Text>
          <Input placeholder="Email" />
          <Input
            placeholder="Password"
            type={showPassword ? "text" : "password"}
          />
          <Button onClick={() => setShowPassword(!showPassword)}>
            {showPassword ? <ViewOffIcon /> : <ViewIcon />}
            {showPassword ? "Hide" : "Show"} Password
          </Button>
          <Button colorScheme="blue">Login</Button>
          <Link href="/signUp">Don't have an account? Sign up</Link>
        </VStack>
      </Center>
    </Container>
  );
}
