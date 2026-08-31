import { ArrowLeft, Github } from 'lucide-react'
import { Link, useLocation } from 'react-router-dom'
import styled from 'styled-components'

import SEO from '../../components/common/SEO'

const Page = styled.section`
  display: grid;
  min-height: min(760px, calc(100vh - 72px));
  place-items: center;
  width: 100%;
  padding: clamp(5rem, 12vw, 9rem) 1.25rem;
  background: ${({ theme }) => theme.colors.background};
`

const Content = styled.div`
  width: min(100%, 720px);
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  padding-top: 2rem;
`

const Eyebrow = styled.p`
  margin: 0 0 1rem;
  color: ${({ theme }) => theme.colors.primary};
  font-family: monospace;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
`

const Title = styled.h1`
  margin: 0;
  max-width: 12ch;
  color: ${({ theme }) => theme.colors.text};
  font-size: clamp(3.5rem, 12vw, 7rem);
  line-height: 0.92;
  letter-spacing: -0.055em;
`

const Description = styled.p`
  max-width: 52ch;
  margin: 1.75rem 0 0;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: clamp(1rem, 2vw, 1.2rem);
  line-height: 1.7;
`

const Actions = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 2rem;
`

const PrimaryLink = styled(Link)`
  display: inline-flex;
  min-height: 48px;
  align-items: center;
  gap: 0.55rem;
  padding: 0.75rem 1rem;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  font-weight: 700;
  text-decoration: none;

  &:hover {
    color: white;
  }
`

const SecondaryLink = styled.a`
  display: inline-flex;
  min-height: 48px;
  align-items: center;
  gap: 0.55rem;
  padding: 0.75rem 1rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  color: ${({ theme }) => theme.colors.text};
  font-weight: 700;
  text-decoration: none;

  &:hover {
    border-color: ${({ theme }) => theme.colors.primary};
    color: ${({ theme }) => theme.colors.primary};
  }
`

export default function NotFound() {
  const location = useLocation()

  return (
    <>
      <SEO
        title="Page not found | Jaden Razo"
        description="The requested page does not exist. Return to Jaden Razo's cloud and DevOps engineering portfolio."
        path={location.pathname}
        noIndex
      />
      <Page aria-labelledby="not-found-title">
        <Content>
          <Eyebrow>HTTP 404 / route not found</Eyebrow>
          <Title id="not-found-title">Wrong route.</Title>
          <Description>
            This path is not part of the current portfolio. Return to the engineering evidence or inspect the source directly on GitHub.
          </Description>
          <Actions>
            <PrimaryLink to="/">
              <ArrowLeft size={18} aria-hidden="true" />
              Return home
            </PrimaryLink>
            <SecondaryLink href="https://github.com/JadenRazo" target="_blank" rel="noopener noreferrer">
              <Github size={18} aria-hidden="true" />
              View GitHub
            </SecondaryLink>
          </Actions>
        </Content>
      </Page>
    </>
  )
}
