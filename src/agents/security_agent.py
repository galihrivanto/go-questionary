from .base_agent import BaseAgent

class SecurityAgent(BaseAgent):
    def get_prompt_template(self) -> str:
        return """
        Perform security analysis for the token:
        {token_data}
        
        Check for:
        1. Contract security patterns
        2. Ownership concentration
        3. Known vulnerabilities
        4. Suspicious transaction patterns
        
        Provide a security risk assessment.
        """
    
    async def analyze(self, token_data: dict) -> dict:
        analysis = await self.chain.arun(token_data=token_data)
        return {
            "agent_type": "security",
            "analysis": analysis,
            "security_score": self._calculate_security_score(token_data)
        }